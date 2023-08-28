package controllers

import (
	"api/initializers"
	"api/models"
	"api/wireguard"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type systemCreateBody struct {
	Name        string `binding:"required"`
	StartedDate string
}

// SystemCreate godoc
// @Summary Create a new system
// @Description Create a new system with the specified name.
// @Tags Systems
// @Accept json
// @Produce json
// @Param body body systemCreateBody true "System details"
// @Success 200 {object} gin.H "System created"
// @Failure 400 {object} gin.H "Fields to read body"
// @Router /systems [post]
func SystemCreate(c *gin.Context) {
	var body systemCreateBody
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "Fields to read body"})
		return
	}

	initializers.DB.Create(&models.System{
		Name:        body.Name,
		StartedDate: body.StartedDate,
		PublicKey:   "sample_Data",
	})

	c.JSON(200, gin.H{"message": "System created"})
}

// SystemShow godoc
// @Summary System Show
// @Description description for SystemShow.
// @Tags Systems
// @Param name path string true "System name"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Produce json
// @Success 200 {object} models.SystemInfo "System information"
// @Failure 400 {object} gin.H "Invalid page number" "Invalid per page number" "Invalid system fetching"
// @Router /systems/{name} [get]
func SystemShow(c *gin.Context) {
	name := c.Param("name")
	page := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "10")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}
	perPageNum, err := strconv.Atoi(perPage)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid per page number"})
		return
	}
	startIdx := (pageNum - 1) * perPageNum
	endIdx := startIdx + perPageNum

	fmt.Print(endIdx)

	var system models.System
	initializers.DB.Model(&models.System{}).Preload("Peers").Where("name = ?", name).First(&system)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid system fetching"})
		return
	}
	var systemInfo models.SystemInfo
	systemInfo.Name = system.Name
	systemInfo.StartedDate = system.StartedDate
	systemInfo.TotalUsage = system.TotalUsage
	selectedPeers := system.Peers
	peersInfo := make([]models.PeerInfo, len(selectedPeers))
	for i, peer := range selectedPeers {
		peersInfo[i].Name = peer.Name
		peersInfo[i].Usage = peer.Usage
		peersInfo[i].DataLimit = peer.DataLimit
		peersInfo[i].BuyDate = peer.BuyDate
		peersInfo[i].ExpireDate = peer.ExpireDate
		peersInfo[i].IsActive = peer.IsActive
	}
	systemInfo.Peers = peersInfo

	c.JSON(200, systemInfo)

}

// SystemShowBasedOnUsage godoc
// @Summary Get systems based on usage
// @Description Retrieve systems based on usage, sorted by usage.
// @Tags Systems
// @Param name path string true "System name"
// @Param page query int false "Page number" default(1)
// @Param per_page query int false "Items per page" default(10)
// @Produce json
// @Success 200 {object} models.SystemInfo "System information"
// @Failure 400 {object} gin.H "Invalid page number" "Invalid per page number"
// @Failure 404 {object} gin.H "System not found"
// @Router /systems/{name}/usage [get]
func SystemShowBasedOnUsage(c *gin.Context) {
	name := c.Param("name")
	page := c.DefaultQuery("page", "1")
	perPage := c.DefaultQuery("per_page", "10")
	pageNum, err := strconv.Atoi(page)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid page number"})
		return
	}
	perPageNum, err := strconv.Atoi(perPage)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid per page number"})
		return
	}
	startIdx := (pageNum - 1) * perPageNum
	endIdx := startIdx + perPageNum
	if endIdx > len(wireguard.SortedPeerBasedOnUsage) {
		endIdx = len(wireguard.SortedPeerBasedOnUsage)
	}

	var system models.System
	initializers.DB.Where("name = ?", name).Preload("Peers").First(&system)
	if system.ID == 0 {
		c.JSON(404, gin.H{"error": "System not found"})
		return
	}
	var result models.SystemInfo
	result.Name = system.Name
	result.StartedDate = system.StartedDate
	result.TotalUsage = system.TotalUsage
	selectedPeers := wireguard.SortedPeerBasedOnUsage[startIdx:endIdx]
	peersInfo := make([]models.PeerInfo, len(selectedPeers))
	for i, peer := range selectedPeers {
		peersInfo[i].Name = peer.Name
		peersInfo[i].Usage = peer.Usage
		peersInfo[i].DataLimit = peer.DataLimit
		peersInfo[i].BuyDate = peer.BuyDate
		peersInfo[i].ExpireDate = peer.ExpireDate
		peersInfo[i].IsActive = peer.IsActive
	}

	result.Peers = peersInfo
	c.JSON(200, result)
}

type systemCreatePeerBody struct {
	Name           string  `json:"name"`
	Phone          string  `json:"phone"`
	Email          string  `json:"email"`
	AllowedIP      string  `json:"allowed_ip"`
	ConfigEndPoint string  `json:"config_end_point"`
	DataLimit      float32 `json:"data_limit"`
	BuyDate        string  `json:"buy_date"`
	ExpireDate     string  `json:"expire_date"`
}

// SystemCreatePeer godoc
// @Summary Create a new peer for a system
// @Description Create a new peer for a system by name.
// @Tags Systems
// @Accept json
// @Produce json
// @Param name path string true "System name"
// @Param body body systemCreatePeerBody true "Peer details"
// @Success 200 {object} gin.H "Peer created"
// @Failure 400 {object} gin.H "Fields to read body"
// @Failure 404 {object} gin.H "System not found"
// @Router /systems/{name}/peers [post]
func SystemCreatePeer(c *gin.Context) {
	name := c.Param("name")
	var body systemCreatePeerBody
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "Fields to read body"})
		return
	}
	var system models.System
	initializers.DB.Where("name = ?", name).Preload("Peer").First(&system)
	if system.ID == 0 {
		c.JSON(404, gin.H{"error": "System not found"})
		return
	}
	publicKey, privateKey := wireguard.GenerateKeys()
	psk := wireguard.GeneratePSK()
	peer := models.Peer{
		Name:           body.Name,
		Phone:          body.Phone,
		Email:          body.Email,
		AllowedIP:      body.AllowedIP,
		ConfigEndPoint: body.ConfigEndPoint,
		PublicKey:      publicKey,
		PrivateKey:     privateKey,
		PreSharedKey:   psk,
		DataLimit:      body.DataLimit,
		BuyDate:        body.BuyDate,
		ExpireDate:     body.ExpireDate,
		SystemID:       system.ID,
	}
	result := initializers.DB.Create(&peer)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Peer not created"})
		return
	}
	c.JSON(200, gin.H{"message": "Peer created"})
}

// SystemReload godoc
// @Summary Reload a system
// @Description Reload a system by name.
// @Tags Systems
// @Param name path string true "System name"
// @Produce json
// @Success 200 {object} gin.H "System reloaded"
// @Failure 404 {object} gin.H "System not found"
// @Router /systems/{name}/reload [post]
func SystemReload(c *gin.Context) {
	name := c.Param("name")
	var system models.System
	initializers.DB.Where("name = ?", name).First(&system)
	if system.ID == 0 {
		c.JSON(404, gin.H{"error": "System not found"})
		return
	}
	err := wireguard.ReloadSystem(system.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Reloading system failed", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "System reloaded"})
}

// TestSystemReload godoc
// @Summary Reload a system
// @Description Reload a system by name.
// @Tags Test
// @Tags Systems
// @Param name path string true "System name"
// @Produce json
// @Success 200 {object} gin.H "System reloaded"
// @Failure 404 {object} gin.H "System not found"
// @Router /test/systems/{name}/reload [post]
func TestSystemReload(c *gin.Context) {
	name := c.Param("name")
	var system models.System
	initializers.DB.Where("name = ?", name).First(&system)
	if system.ID == 0 {
		c.JSON(404, gin.H{"error": "System not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "System reloaded"})
}

// TestSystemCreatePeer godoc
// @Summary Create a new peer for a system
// @Description Create a new peer for a system by name.
// @Tags Systems
// @Tags Test
// @Accept json
// @Produce json
// @Param name path string true "System name"
// @Param body body systemCreatePeerBody true "Peer details"
// @Success 200 {object} gin.H "Peer created"
// @Failure 400 {object} gin.H "Fields to read body"
// @Failure 404 {object} gin.H "System not found"
// @Router /test/systems/{name}/peers [post]
func TestSystemCreatePeer(c *gin.Context) {
	name := c.Param("name")
	var body systemCreatePeerBody
	if c.Bind(&body) != nil {
		c.JSON(400, gin.H{"error": "Fields to read body"})
		return
	}
	var system models.System
	initializers.DB.Where("name = ?", name).Preload("Peer").First(&system)
	if system.ID == 0 {
		c.JSON(404, gin.H{"error": "System not found"})
		return
	}
	publicKey, privateKey := "sample_Data", "sample_Data"
	psk := "sample_Data"
	peer := models.Peer{
		Name:           body.Name,
		Phone:          body.Phone,
		Email:          body.Email,
		AllowedIP:      body.AllowedIP,
		ConfigEndPoint: body.ConfigEndPoint,
		PublicKey:      publicKey,
		PrivateKey:     privateKey,
		PreSharedKey:   psk,
		DataLimit:      body.DataLimit,
		BuyDate:        body.BuyDate,
		ExpireDate:     body.ExpireDate,
		SystemID:       system.ID,
	}
	result := initializers.DB.Create(&peer)
	if result.Error != nil {
		c.JSON(400, gin.H{"error": "Peer not created"})
		return
	}
	c.JSON(200, gin.H{"message": "Peer created"})
}
