package data

import (
	"regexp"
	"strings"
)

// TagsProcessor Helper class to extract and process search tags from objects.
// The search tags can be kept individually or embedded as hash tags inside text like
// "This text has #hash_tag that can be used for search."
var TagsProcessor = &_TTagsProcessor{}

type _TTagsProcessor struct{}

var normalizeTagRegex = regexp.MustCompile("(_|#)+")
var compressTagRegex = regexp.MustCompile("( |_|#)")
var splitTagRegex = regexp.MustCompile("(,|;)+")
var hashTagRegex = regexp.MustCompile("#\\w+")

// NormalizeTag normalizes a tag by replacing special symbols like '_' and '#' with spaces.
// When tags are normalized then can be presented to user in similar shape and form.
//	Parameters: tag string the tag to normalize.
//	Returns: string a normalized tag.
func (c *_TTagsProcessor) NormalizeTag(tag string) string {
	tag = normalizeTagRegex.ReplaceAllString(tag, " ")
	return strings.Trim(tag, " \t\r\n")
}

// CompressTag compress a tag by removing special symbols like spaces, '_' and '#'
// and converting the tag to lower case. When tags are compressed they can be matched in search queries.
//	Parameters: tag string the tag to compress.
//	Returns: string a compressed tag.
func (c *_TTagsProcessor) CompressTag(tag string) string {
	tag = compressTagRegex.ReplaceAllString(tag, "")
	return strings.ToLower(tag)
}

// EqualTags compares two tags using their compressed form.
//	Parameters:
//		- tag1 string the first tag.
//		- tag2 string the second tag.
//	Returns: bool true if the tags are equal and false otherwise.
func (c *_TTagsProcessor) EqualTags(tag1 string, tag2 string) bool {
	if tag1 == "" && tag2 == "" {
		return true
	}
	if tag1 == "" || tag2 == "" {
		return false
	}
	return strings.Compare(c.CompressTag(tag1), c.CompressTag(tag2)) == 0
}

// NormalizeTags normalizes a list of tags.
//	Parameters: tags []string the tags to normalize.
//	Returns: []string a list with normalized tags.
func (c *_TTagsProcessor) NormalizeTags(tags []string) []string {
	for index := 0; index < len(tags); index++ {
		tags[index] = c.NormalizeTag(tags[index])
	}
	return tags
}

// NormalizeTagList normalizes a comma-separated list of tags.
//	Parameters: tagList string a comma-separated list of tags to normalize.
//	Returns: []string a list with normalized tags.
func (c *_TTagsProcessor) NormalizeTagList(tagList string) []string {
	tags := splitTagRegex.Split(tagList, -1)
	return c.NormalizeTags(tags)
}

// CompressTags compresses a list of tags.
//	Parameters: tags []string the tags to compress.
//	Returns: []string a list with normalized tags.
func (c *_TTagsProcessor) CompressTags(tags []string) []string {
	for index := 0; index < len(tags); index++ {
		tags[index] = c.CompressTag(tags[index])
	}
	return tags
}

// CompressTagList compresses a comma-separated list of tags.
//	Parameters: tagList string a comma-separated list of tags to compress.
//	Returns: []string a list with compressed tags.
func (c *_TTagsProcessor) CompressTagList(tagList string) []string {
	tags := splitTagRegex.Split(tagList, -1)
	return c.CompressTags(tags)
}

// ExtractHashTags extracts hash tags from a text.
//	Parameters: text string a text that contains hash tags
//	Returns: []string a list with extracted and compressed tags.
func (c *_TTagsProcessor) ExtractHashTags(text string) []string {
	if text == "" {
		return []string{}
	}

	hashTags := hashTagRegex.FindAllString(text, -1)
	tags := make([]string, 0, len(hashTags))

	// Filter duplicates
	for index := 0; index < len(hashTags); index++ {
		tag := c.CompressTag(hashTags[index])
		duplicate := false
		for index2 := 0; index2 < len(tags); index2++ {
			if strings.Compare(tags[index2], tag) == 0 {
				duplicate = false
				break
			}
		}
		if !duplicate {
			tags = append(tags, tag)
		}
	}

	return tags
}
