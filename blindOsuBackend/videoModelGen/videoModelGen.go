// Shit dont work (.wav file is wrong format)
package videomodelgen

import (
	"fmt"
	"os"

	"github.com/mjibson/go-dsp/fft"
	wav "github.com/mjibson/go-dsp/wav"
	"github.com/mjibson/go-dsp/window"
	"github.com/mrnikho/yingo"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"github.com/vidur2/blindOsuBackend/types"
)

const DURATION = 0.015

func GenerateCoordPoints() ([]types.RelativeModelCoord, error) {

	defer os.Remove("../audio.wav")

	f, err := os.Open("../audio.wav")

	if err != nil {
		return nil, err
	}

	wav, err := wav.New(f)

	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	var modelCoords []types.RelativeModelCoord

	freqs := getFreq(wav)

	maxVal := calculateMax(freqs)

	modelCoords = make([]types.RelativeModelCoord, len(freqs))

	for idx, freq := range freqs {
		modelCoords[idx] = freq.ConvertAbsModelCoord(maxVal)
	}

	return modelCoords, nil
}

func yingoUse() ([]float32, error) {
	f, err := os.Open("../audio.wav")

	if err != nil {
		return nil, err
	}

	yin := yingo.Yin{}
	pcm, err := wav.New(f)

	if err != nil {
		return nil, err
	}

	yin.YinInit(100, 1)
	samplesInt, _ := pcm.ReadSamples(pcm.Samples)
	samplesCasted := samplesInt.([]int16)
	// fmt.Println(thingFinal)

	pitches := make([]float32, 0)
	i := 0
	for i < pcm.Samples-100 {
		slice := samplesCasted[i : i+100]
		i += 100
		pitch := yin.GetPitch(&slice)
		pitches = append(pitches, pitch)
	}
	return pitches, nil
}

func getFreq(wav *wav.Wav) []types.AbsModelCoord {
	fmt.Println(wav.Duration.Seconds())
	wavUnparsed, _ := wav.ReadSamples(32768)
	wavParsed := wavUnparsed.([]int16)

	sliced := make([]float64, len(wavParsed))

	for idx, val := range wavParsed {
		sliced[idx] = float64(val)
	}

	window.Apply(sliced, window.Hamming)

	w := fft.FFTReal(sliced)
	fmt.Println(len(w))

	retVal := make([]types.AbsModelCoord, len(w))

	for idx, coord := range w {
		retVal[idx] = types.AbsModelCoord{
			AbsolutePitch: coord,
			Time:          float64(idx) * (wav.Duration.Seconds() / 32768),
		}
	}

	return retVal

}

func calculateMax(absPitch []types.AbsModelCoord) complex128 {
	maxReal := 0.
	maxImag := 0.

	for _, val := range absPitch {
		if real(val.AbsolutePitch) > maxReal {
			maxReal = real(val.AbsolutePitch)
		}

		if imag(val.AbsolutePitch) > maxImag {
			maxImag = imag(val.AbsolutePitch)
		}
	}

	return complex(maxReal, maxImag)
}

func TransposeMp3File() error {
	// os.Remove("../audio.wav")
	err := ffmpeg_go.Input("/Users/vidurmodgil/Desktop/ProgrammingProjects/blindOsuFull/blindOsuBackend/audio.mp3").Output("../audio.wav").Run()
	// cmd := exec.Command("ffmpeg -i /Users/vidurmodgil/Desktop/ProgrammingProjects/blindOsuFull/blindOsuBackend/audio.mp3 /Users/vidurmodgil/Desktop/ProgrammingProjects/blindOsuFull/blindOsuBackend/audio.wav")

	// err := cmd.Run()

	return err
}
