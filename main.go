package main

import (
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/polly"
)

// ...
func main() {
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	sessionToken := os.Getenv("AWS_SESSION_TOKEN")
	region := "us-west-2"

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, sessionToken),
	})
	if err != nil {
		fmt.Println("Failed to create AWS session:", err)
		return
	}
	// Initialize Polly service client
	pollySvc := polly.New(sess)

	// Text to be converted to speech
	text := "The Plant Export Operations Manual (PEOM) and associated Work Instructions for the export of grain and plant products provide procedural references for the sampling, inspection and certification of export grain and plant products across Australia. Both PEOM and Work Instructions should be read in conjunction with the Export Control (Plants and Plant Products) Order 2011."

	// Specify the voice you want to use
	voiceID := "Joanna"

	// Generate the input parameters for the SynthesizeSpeech operation
	input := &polly.SynthesizeSpeechInput{
		OutputFormat: aws.String("mp3"),
		Text:         aws.String(text),
		VoiceId:      aws.String(voiceID),
	}
	output, err := pollySvc.SynthesizeSpeech(input)
	if err != nil {
		fmt.Println("Failed to synthesize speech:", err)
		return
	}
	audioFileName := "output.mp3"
	audioFile, err := os.Create(audioFileName)
	if err != nil {
		fmt.Println("Failed to create audio file:", err)
		return
	}
	defer audioFile.Close()

	bytes, err := io.ReadAll(output.AudioStream)

	if err != nil {
		fmt.Println("Error reading data:", err)
		return
	}

	_, err = audioFile.Write(bytes)
	if err != nil {
		fmt.Println("Failed to write audio data to file:", err)
		return
	}
	fmt.Printf("Speech synthesized and saved to %s\n", audioFileName)
}
