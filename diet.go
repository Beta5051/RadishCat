package RadishCat

import (
	"github.com/Beta5051/NeisGo"
	"github.com/bwmarrin/discordgo"
	"strings"
	"time"
)

var areaCodeMap = map[string]string{
	"서울": NeisGo.SEOUL,
	"부산": NeisGo.BUSAN,
	"대그": NeisGo.DAEGU,
	"인천": NeisGo.INCHEON,
	"대전": NeisGo.DAEJOEN,
	"울산": NeisGo.ULSAN,
	"세종": NeisGo.SEJONG,
	"경기": NeisGo.GYEONGII,
	"강원": NeisGo.GANGWON,
	"충북": NeisGo.CHUNGCHEONG_NORTH,
	"충남": NeisGo.CHUNGCHEONG_SOUTH,
	"전북": NeisGo.JEOLLA_NORTH,
	"전남": NeisGo.JEOLLA_SOUTH,
	"경북": NeisGo.GYEONGSANG_NORTH,
	"경남": NeisGo.GYEONGSANG_SOUTH,
	"제주": NeisGo.JEJU,
	"재외": NeisGo.OVERSEAS,
}

type dietCommand struct {
	Area string
	School string
	Kind string
	Date string
}

func (d *dietCommand) Run(rc *RadishCat, m *discordgo.Message) error {
	if areaCode, ok := areaCodeMap[d.Area]; ok{
		schoolInfos, err := rc.Neis.GetSchoolInfo(NeisGo.SchoolInfoFactor{
			ATPT_OFCDC_SC_CODE: areaCode,
			SCHUL_NM: d.School,
			SCHUL_KND_SC_NM: d.Kind,
		})
		if err != nil {
			return err
		}

		if d.Date == "" {
			d.Date = time.Now().Format("20060102")
		}

		dietInfos, err := rc.Neis.GetDietInfo(NeisGo.DietInfoFactor{
			ATPT_OFCDC_SC_CODE: areaCode,
			SD_SCHUL_CODE: schoolInfos[0].SD_SCHUL_CODE,
			MLSV_YMD: d.Date,
		})
		if err != nil {
			return err
		}

		embed := rc.NewEmbed(dietInfos[0].SCHUL_NM + " 급식")
		embed.SetDescription("날짜: " + d.Date)
		for _, dietInfo := range dietInfos {
			embed.AddField(dietInfo.MMEAL_SC_NM, strings.ReplaceAll(dietInfo.DDISH_NM, "<br/>", "\n"), false)
		}

		_, err = rc.SendMessage(m.ChannelID, embed)
		return err
	}
	return nil
}