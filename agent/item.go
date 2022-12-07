package agent

import (
	"github.com/spf13/viper"
	"log"
	"math/rand"
	"mouban/consts"
	"mouban/dao"
	"mouban/util"
	"strconv"
	"time"
)

func runItem(index int) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r, "item agent", index, "crashed  => ", util.GetCurrentGoroutineStack())
		}
		time.Sleep(time.Second * 1)
	}()

	types := []consts.Type{consts.TypeBook, consts.TypeMovie, consts.TypeGame}
	rand.Shuffle(len(types), func(i, j int) { types[i], types[j] = types[j], types[i] })

	for _, t := range types {
		schedule := dao.SearchScheduleByStatus(t.Code, consts.ScheduleStatusToCrawl)
		if schedule != nil {
			changed := dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawling, consts.ScheduleStatusToCrawl)
			if changed {
				log.Println("item", index, "start", t.Name, strconv.FormatUint(schedule.DoubanId, 10))
				processItem(t.Code, schedule.DoubanId)
				dao.CasScheduleStatus(schedule.DoubanId, schedule.Type, consts.ScheduleStatusCrawled, consts.ScheduleStatusCrawling)
				log.Println("item", index, "end", t.Name, strconv.FormatUint(schedule.DoubanId, 10))
			}
			break
		}
	}

}

func init() {
	if !viper.GetBool("agent.enable") {
		log.Println("item agent disabled")
		return
	}
	concurrency := viper.GetInt("agent.item.concurrency")
	for i := 0; i < concurrency; i++ {
		j := i + 1
		go func() {
			for {
				runItem(j)
			}
		}()
	}

	log.Println(concurrency, " item agent(s) enabled")
}
