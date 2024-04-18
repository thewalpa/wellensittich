package config

type WellensittichConfig struct {
	DevServer   string `json:"dev_server"`
	Token       string `json:"token"`
	WhisperHost string `json:"whisper_asr_webservice"`
	TenorKey    string `json:"tenor_key"`
}
