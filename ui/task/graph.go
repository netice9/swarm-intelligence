package task

import (
	"fmt"
	"strings"
	"time"

	"github.com/netice9/swarm-intelligence/model/stats"
	"gitlab.netice9.com/dragan/go-reactor/core"
)

var graphUI = core.MustParseDisplayModel(`
  <div>
	  <bs.Panel id="goal_panel" header="CPU Stats">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 450 130" width="100%" class="chart">
				<g transform="translate(10,20)">
					<path d="M28 0h3M28 100h3M31 100v3" strokeWidth="1px" stroke="#333"/>
					<path d="M31 0v100M31 100h400" strokeWidth="1px" stroke="#333"/>
					<polyline transform="translate(32,0)" id="cpuLine" fill="none" stroke="#0074d9" strokeWidth="1" points=""/>
					<g fontSize="8px" fontFamily="Georgia" fill="#333">
						<g textAnchor="end">
							<text id="maxCPU" x="26" y="2">100 %</text>
							<text x="26" y="102">0 %</text>
						</g>
					</g>
				</g>
			</svg>
		</bs.Panel>
		<bs.Panel id="goal_panel" header="Memory Stats">
			<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 450 130" width="100%" class="chart">
				<g transform="translate(10,20)">
				  <polyline transform="translate(32,0)" id="memLine" fill="none" stroke="#0074d9" strokeWidth="1" points=""/>
					<path d="M28 0h3M28 100h3M31 100v3" strokeWidth="1px" stroke="#333"/>
					<path d="M31 0v100M31 100h400" strokeWidth="1px" stroke="#333"/>
					<g fontSize="8px" fontFamily="Georgia" fill="#333">
						<g textAnchor="end">
							<text id="maxMem" x="26" y="2">100 MB</text>
							<text x="26" y="102">0 MB</text>
						</g>
					</g>
				</g>
			</svg>
	  </bs.Panel>
	</div>
`)

type sample struct {
	time  time.Time
	value float64
}

func renderGraph(entries []stats.Entry) *core.DisplayModel {
	cpuSamples := []sample{}
	memSamples := []sample{}
	for _, e := range entries {
		cpuSamples = append(cpuSamples, sample{e.Time, float64(e.CPU) / 1e7})
		memSamples = append(cpuSamples, sample{e.Time, float64(e.Memory) / (1024 * 1024)})
	}

	g := graphUI.DeepCopy()

	cpuPoints, maxCPU := timeSeriesToLines(cpuSamples, 400, 100, 0.1)
	g.SetElementAttribute("cpuLine", "points", cpuPoints)
	g.SetElementText("maxCPU", fmt.Sprintf("%.1f%%", maxCPU))

	memPoints, maxMem := timeSeriesToLines(memSamples, 400, 100, 0.4)
	g.SetElementAttribute("memLine", "points", memPoints)
	g.SetElementText("maxMem", fmt.Sprintf("%.1f MB", maxMem))

	return g
}

func timeSeriesToLines(samples []sample, width, height int, lowestMax float64) (string, float64) {
	if len(samples) == 0 {
		return "", 0
	}
	// sample := samples[0]

	minValue := samples[0].value
	maxValue := lowestMax
	minTime := samples[0].time
	maxTime := samples[0].time

	for _, sample := range samples {
		if minValue > sample.value {
			minValue = sample.value
		}
		if maxValue < sample.value {
			maxValue = sample.value
		}
		if minTime.After(sample.time) {
			minTime = sample.time
		}
		if maxTime.Before(sample.time) {
			maxTime = sample.time
		}
	}

	points := []string{}

	for _, sample := range samples {
		normalisedTime := float64(sample.time.UnixNano()-minTime.UnixNano()) / float64((time.Second * 120).Nanoseconds())

		scaledTime := int(normalisedTime * float64(width))

		normalisedValue := 1.0 - float64(sample.value-minValue)/float64(maxValue-minValue)
		scaledValue := int(normalisedValue * float64(height))
		points = append(points, fmt.Sprintf("%d,%d", scaledTime, scaledValue))
	}

	return strings.Join(points, " "), maxValue

}
