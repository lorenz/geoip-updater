package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jasonlvhit/gocron"
)

type MMUpdater struct {
	EditionID string
	UserID    uint64
	Storage   Storage
}

func NewFSMMUpdater(editionID string, userID uint64, path string) MMUpdater {
	return MMUpdater{
		EditionID: editionID,
		UserID:    userID,
		Storage:   NewFileSystemStorage(path),
	}
}

func (m MMUpdater) Update() (bool, error) {
	md5 := m.Storage.GetMD5(m.EditionID)
	reader, err := UpdateReader(m.EditionID, m.UserID, md5)
	if err != nil {
		return false, err
	}
	if reader != nil {
		if err = m.Storage.Update(m.EditionID, reader); err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func UpdateJob(m MMUpdater) {
	log.Printf("Starting update of %v", m.EditionID)
	hasUpdated, err := m.Update()
	if err != nil {
		log.Printf("Update of %v failed: %v", m.EditionID, err)
		return
	}
	if hasUpdated {
		log.Printf("%v has been updated to the newest version", m.EditionID)
	} else {
		log.Printf("%v has no new version available", m.EditionID)
	}
}

func main() {
	editionIDs := strings.Split(os.Getenv("EDITION_IDS"), ",")
	rawUserID := os.Getenv("USER_ID")

	userID, err := strconv.ParseUint(rawUserID, 10, 64)
	if err != nil {
		log.Panicf("USER_ID is not a valid number")
	}

	if len(editionIDs) == 0 {
		log.Panicf("No EDITION_IDs given, please specifiy at least one")
	}

	log.Printf("Started geoipd-updater")

	for _, editionID := range editionIDs {
		mmupdater := NewFSMMUpdater(editionID, userID, "/data")
		go UpdateJob(mmupdater)
		gocron.Every(1).Hour().Do(UpdateJob, mmupdater)
	}
	<-gocron.Start()
}
