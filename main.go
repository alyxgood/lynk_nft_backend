package main

import (
	"alyx_nft_backend/consts"
	"alyx_nft_backend/models"
	"alyx_nft_backend/utils"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickmn/go-cache"
	"github.com/urfave/cli"
)

func main() {
	if os.Getenv("debugPProf") == "true" {
		go func() {
			log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
		}()
	}

	app := cli.NewApp()

	app.Name = "alyx_nft_backend"
	app.Version = "v0.1.0"
	app.Description = "alyx_nft_backend"
	server := NewService()
	app.Action = server.Start

	_ = app.Run(os.Args)
}

type Service struct {
	mCache     *cache.Cache
	mCacheTime time.Duration
}

func NewService() *Service {
	return &Service{}
}

func (svc *Service) Start(ctx *cli.Context) error {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("got panic: %v\n", err)
		}
	}()

	svc.initMainLogger()

	svc.mCacheTime = time.Duration(1) * time.Minute
	svc.mCache = cache.New(svc.mCacheTime, svc.mCacheTime)

	router := gin.Default()
	router.Use(utils.Cors())
	router.GET(path.Join("/alyx", "/nft/:tokenId"), svc.httpQueryNFTInfo)

	return router.Run(fmt.Sprintf("0.0.0.0:%d", consts.Port))
}

func (svc *Service) initMainLogger() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
}

func (svc *Service) nftFunc(tokenId int) ([]interface{}, bool, error) {
	const dataExistsPrefix = "0x4f558e79"
	dataExists := fmt.Sprintf("%s%s", dataExistsPrefix, fmt.Sprintf("%064x", tokenId))

	resExists, err := utils.QueryNFTInfo(consts.LYNKNFTAddress, dataExists, consts.JsonRpc)
	if err != nil {
		log.Println("QueryNFTInfo failed.")
		return nil, false, err
	}
	var exist []interface{}
	if resExists.Result != "" {
		var outputParameters []string
		outputParameters = append(outputParameters, "bool")
		exist, err = utils.Decode(outputParameters, strings.TrimPrefix(resExists.Result, "0x"))
		if err != nil {
			log.Println("Decode failed.")
			return nil, false, err
		}
	}

	var arr []interface{}
	if len(exist) != 0 {
		if exist[0].(bool) {
			const dataNFTPrefix = "0xdcefcebc"
			dataNFT := fmt.Sprintf("%s%s", dataNFTPrefix, fmt.Sprintf("%064x", tokenId))

			res, err := utils.QueryNFTInfo(consts.LYNKNFTAddress, dataNFT, consts.JsonRpc)
			if err != nil {
				log.Println("QueryNFTInfo failed.")
				return nil, false, err
			}

			if res.Result != "" {
				var outputParameters []string
				outputParameters = append(outputParameters, "uint256[]")
				arr, err = utils.Decode(outputParameters, strings.TrimPrefix(res.Result, "0x"))
				if err != nil {
					log.Println("Decode failed.")
					return nil, false, err
				}
			}
		} else {
			return nil, false, nil
		}
	}

	return arr, true, nil
}

func (svc *Service) httpQueryNFTInfo(c *gin.Context) {
	tokenId := c.Param("tokenId")
	tokenIdParam, err := strconv.Atoi(tokenId)
	if err != nil {
		log.Printf("BadRequest.: %v\n", err)
		c.JSON(http.StatusBadRequest, models.Response{
			Code:    400,
			Message: "Something went wrong, Please try again later.",
		})
		return
	}
	var charismaStr, vitalityStr, intellectStr, dexterityStr string

	keyCharisma := fmt.Sprintf("%s-%s", tokenId, "charisma")
	keyVitality := fmt.Sprintf("%s-%s", tokenId, "vitality")
	keyIntellect := fmt.Sprintf("%s-%s", tokenId, "intellect")
	keyDexterity := fmt.Sprintf("%s-%s", tokenId, "dexterity")
	cacheCharisma, ok1 := svc.mCache.Get(keyCharisma)
	cacheVitality, ok2 := svc.mCache.Get(keyVitality)
	cacheIntellect, ok3 := svc.mCache.Get(keyIntellect)
	cacheDexterity, ok4 := svc.mCache.Get(keyDexterity)
	if !ok1 || !ok2 || !ok3 || !ok4 {
		var arr []interface{}
		var exist bool
		arr, exist, err = svc.nftFunc(tokenIdParam)
		if err != nil || len(arr) == 0 {
			log.Printf("nftTask failed.: %v\n", err)
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "Something went wrong, Please try again later.",
			})
			return
		}
		if !exist {
			c.JSON(http.StatusInternalServerError, models.Response{
				Code:    http.StatusInternalServerError,
				Message: "nft not exist!",
			})
			return
		}

		arrBigInt := arr[0].([]*big.Int)
		charismaStr = arrBigInt[0].String()
		vitalityStr = arrBigInt[1].String()
		intellectStr = arrBigInt[2].String()
		dexterityStr = arrBigInt[3].String()

		svc.mCache.Set(keyCharisma, charismaStr, svc.mCacheTime)
		svc.mCache.Set(keyVitality, vitalityStr, svc.mCacheTime)
		svc.mCache.Set(keyIntellect, intellectStr, svc.mCacheTime)
		svc.mCache.Set(keyDexterity, dexterityStr, svc.mCacheTime)

	} else {
		charismaStr = cacheCharisma.(string)
		vitalityStr = cacheVitality.(string)
		intellectStr = cacheIntellect.(string)
		dexterityStr = cacheDexterity.(string)
	}

	var strArr = [4]string{"charisma", "vitality", "intellect", "dexterity"}
	var valueArr = [4]string{charismaStr, vitalityStr, intellectStr, dexterityStr}
	var attribute = make([]models.Attribute, 0)
	for index := range strArr {
		attribute = append(attribute, models.Attribute{
			TraitType: strArr[index],
			Value:     valueArr[index],
		})
	}

	img := fmt.Sprintf(consts.TokenImage, tokenIdParam)

	c.JSON(http.StatusOK, models.ResNFT{
		Description: "LYNK NFT",
		Image:       img,
		Name:        "LYNK NFT",
		Attributes:  attribute,
	})

}
