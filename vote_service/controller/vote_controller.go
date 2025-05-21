package controller
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	"example.com/se/config"
	"example.com/se/entity"
)

// Struct สำหรับ response จากแต่ละ service
type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Candidate struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Election struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

// ฟังก์ชันดึงข้อมูล User จาก user_service
func fetchUser(userID uint) (*User, error) {
	url := fmt.Sprintf("http://user_service:8001/api/user/%d", userID)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to fetch user:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("user service returned status %d", resp.StatusCode)
	}

	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

// ฟังก์ชันดึงข้อมูล Candidate จาก candidate_service
func fetchCandidate(candidateID uint) (*Candidate, error) {
	url := fmt.Sprintf("http://candidate_service:8003/api/candidate/%d", candidateID)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to fetch candidate:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("candidate service returned status %d", resp.StatusCode)
	}

	var candidate Candidate
	if err := json.NewDecoder(resp.Body).Decode(&candidate); err != nil {
		return nil, err
	}

	return &candidate, nil
}

// ฟังก์ชันดึงข้อมูล Election จาก election_service
func fetchElection(electionID uint) (*Election, error) {
	url := fmt.Sprintf("http://election_service:8002/api/election/%d", electionID)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Failed to fetch election:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("election service returned status %d", resp.StatusCode)
	}

	var election Election
	if err := json.NewDecoder(resp.Body).Decode(&election); err != nil {
		return nil, err
	}

	return &election, nil
}

// ตัวอย่าง handler ดึง vote พร้อมข้อมูล detail จาก service ต่างๆ
func GetVoteWithDetails(c *gin.Context) {
	db := config.DB()
	var votes []entity.Votes

	if err := db.Find(&votes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	type VoteWithDetails struct {
		entity.Votes
		User      *User      `json:"user,omitempty"`
		Candidate *Candidate `json:"candidate,omitempty"`
		Election  *Election  `json:"election,omitempty"`
	}

	var result []VoteWithDetails
	for _, v := range votes {
		user, _ := fetchUser(v.UserID)
		candidate, _ := fetchCandidate(v.CandidateID)
		election, _ := fetchElection(v.ElectionID)

		result = append(result, VoteWithDetails{
			Votes:     v,
			User:      user,
			Candidate: candidate,
			Election:  election,
		})
	}

	c.JSON(http.StatusOK, result)
}

// CreateVote รับ JSON มาสร้าง vote ใหม่
func CreateVote(c *gin.Context) {
	var vote entity.Votes

	if err := c.ShouldBindJSON(&vote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// กำหนดเวลาปัจจุบัน
	vote.Timestamp = time.Now()

	db := config.DB()
	if err := db.Create(&vote).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create vote"})
		return
	}

	c.JSON(http.StatusCreated, vote)
}

// GetAllVotes ดึงข้อมูลโหวตทั้งหมด
func GetAllVotes(c *gin.Context) {
	db := config.DB()
	var votes []entity.Votes

	if err := db.Find(&votes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get votes"})
		return
	}

	c.JSON(http.StatusOK, votes)
}

// GetVotesByCandidate ดึงข้อมูลโหวตทั้งหมดตาม candidate_id
func GetVotesByCandidate(c *gin.Context) {
	candidateIDStr := c.Param("candidate_id")
	candidateID, err := strconv.Atoi(candidateIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}

	db := config.DB()
	var votes []entity.Votes

	if err := db.Where("candidate_id = ?", candidateID).Find(&votes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get votes"})
		return
	}

	c.JSON(http.StatusOK, votes)
}

// DeleteVote ลบ vote โดย ID
func DeleteVote(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vote ID"})
		return
	}

	db := config.DB()
	if err := db.Delete(&entity.Votes{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete vote"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote deleted"})
}