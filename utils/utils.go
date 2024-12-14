package utils

import (
	"crypto/rand"
	"discuz/database/models"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"
)

func RenderTemplate(w http.ResponseWriter, page string, data interface{}) {
	path := fmt.Sprintf("cmd/static/%s.html", page)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, "Error in parsing html file", 500)
		return
	}

	tmpl.Execute(w, data)
}

func GenerateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Could not generate Token")
	}
	return base64.URLEncoding.EncodeToString(bytes)
}

func BuildTree(threads []*models.Threads) []*models.Threads {

	treeMap := []*models.Threads{}

	for _, thread := range threads {
		if !thread.Parent_id.Valid {
			treeMap = append(treeMap, thread)
		} else {
			SearchTree(treeMap, thread)
		}
	}

	for _, thread := range treeMap {
		thread.Root_id = thread.ID
		AssignRootThread(thread, thread.ID)
	}

	return treeMap
}

func SearchTree(treeMap []*models.Threads, target *models.Threads) {
	for _, thread := range treeMap {
		if thread.ID == target.Parent_id.String {
			thread.Children = append(thread.Children, target)
			return
		}
		if len(thread.Children) == 0 {
			continue
		} else {
			SearchTree(thread.Children, target)
		}
	}
}

func AssignRootThread(thread *models.Threads, id string) {
	for _, thread := range thread.Children {
		thread.Root_id = id
		if len(thread.Children) != 0 {
			AssignRootThread(thread, id)
		}
	}
}

func ParseDateString(datestring *string) {

	ist, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		fmt.Println("Failed to load IST timezone:", err)
		return
	}

	t, err := time.Parse("2006-01-02T15:04:05Z", *datestring)
	if err != nil {
		return
	}

	t = t.In(ist)

	var minute string

	if t.Minute() < 10 {
		minute = fmt.Sprintf("0%d", t.Minute())
	}

	*datestring = fmt.Sprintf("%d %s %d, %d:%s", t.Day(), t.Month(), t.Year(), t.Hour(), minute)
}
