package libtools

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

func ConvertSticker(client *whatsmeow.Client, v *events.Message) error {
	img := v.Message.GetImageMessage()
	ctx := context.Background()
	if img == nil {
		return fmt.Errorf("tidak ada gambar yang dikirim")
	}

	// Download media
	data, err := client.Download(img)
	if err != nil {
		return fmt.Errorf("gagal download gambar: %w", err)
	}

	// Simpan file sementara
	tempInput := "input.jpg"
	tempOutput := "output.webp"
	if err := os.WriteFile(tempInput, data, 0644); err != nil {
		return fmt.Errorf("gagal simpan input file: %w", err)
	}

	// Konversi ke WebP sesuai spesifikasi WhatsApp Sticker
	cmd := exec.Command("ffmpeg",
		"-i", tempInput,
		"-vf", "scale=512:512:force_original_aspect_ratio=decrease,pad=512:512:(ow-iw)/2:(oh-ih)/2:color=white",
		"-vcodec", "libwebp",
		"-lossless", "0",
		"-q:v", "80",
		"-loop", "0",
		"-preset", "default",
		"-an", "-vsync", "0",
		tempOutput,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("gagal konversi ke webp: %v\n%s", err, output)
	}

	// Baca hasil webp
	webpData, err := os.ReadFile(tempOutput)
	if err != nil {
		return fmt.Errorf("gagal membaca output webp: %w", err)
	}

	// Upload ke server WA - pakai MediaImage
	uploaded, err := client.Upload(ctx, webpData, whatsmeow.MediaImage)
	if err != nil {
		return fmt.Errorf("upload gagal: %w", err)
	}

	// Buat pesan sticker
	stickerMsg := &waProto.Message{
		StickerMessage: &waProto.StickerMessage{
			URL:           &uploaded.URL,
			Mimetype:      proto.String("image/webp"),
			MediaKey:      uploaded.MediaKey,
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    &uploaded.FileLength,
			DirectPath:    &uploaded.DirectPath,
		},
	}

	// Kirim stiker ke pengirim
	_, err = client.SendMessage(ctx, v.Info.Chat, stickerMsg)

	// Bersihkan file sementara
	_ = os.Remove(tempInput)
	_ = os.Remove(tempOutput)

	if err != nil {
		return fmt.Errorf("gagal kirim sticker: %w", err)
	}

	return nil
}
