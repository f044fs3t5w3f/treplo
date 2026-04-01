package salute

type RecognizeChunk struct {
	Text           string `json:"text"`
	NormalizedText string `json:"normalized_text"`
	Start          string `json:"start"`
	End            string `json:"end"`
	WordAlignments []struct {
		Word  string `json:"word"`
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"word_alignments"`
}
