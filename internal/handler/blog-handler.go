package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/russross/blackfriday"
	"github.com/vitezprchal/waitlistgo/internal/models"
	"github.com/vitezprchal/waitlistgo/website/view"
)

const PER_PAGE = 10

type BlogHandler struct {
	blogs []*models.Blog
}

func read_seo(seo_file string) *models.SEO {
	seo, err := os.Open(seo_file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer seo.Close()

	seo_info, err := os.Stat(seo_file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	seo_size := seo_info.Size()
	seo_content := make([]byte, seo_size)
	_, err = seo.Read(seo_content)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var seo_model *models.SEO
	err = json.Unmarshal(seo_content, &seo_model)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return seo_model
}

func read_md(md_file string) template.HTML {
	file, err := ioutil.ReadFile(md_file)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	html := template.HTML(blackfriday.MarkdownCommon([]byte(file)))
	return html
}

func NewBlogHandler() *BlogHandler {
	dir, err := os.Open("./website/blogs")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var new_blogs []*models.Blog
	for _, file_info := range files {
		if file_info.IsDir() {
			seo := read_seo(fmt.Sprintf("./website/blogs/%s/seo.json", file_info.Name()))
			content := read_md(fmt.Sprintf("./website/blogs/%s/content.md", file_info.Name()))

			blog := &models.Blog{
				Slug:    file_info.Name(), //slugify
				SEO:     seo,
				Content: content,
			}

			new_blogs = append(new_blogs, blog)
		}
	}

	// var new_blogs []*models.Blog

	// for i := 0; i < 100; i++ {
	// 	seo := &models.SEO{
	// 		Title:       "Blog",
	// 		Description: "Blog Description",
	// 		Keywords:    "Blog",
	// 		AuthorName:  "Vitezslav",
	// 		ImageURL:    "https://www.something.some/assets/logo-image.png",
	// 		PageURL:     "https://www.something.some",
	// 	}

	// 	content := template.HTML("<h1>Blog Post</h1><p>This is a blog post</p>")
	// 	blog := &models.Blog{
	// 		Slug:    fmt.Sprintf("blog-post-%d", i),
	// 		SEO:     seo,
	// 		Content: content,
	// 	}

	// 	new_blogs = append(new_blogs, blog)
	// }
	// seo := &models.SEO{
	// 	Title:       "Last Blog",
	// 	Description: "This is a blog post",
	// 	Keywords:    "blog, post",
	// 	AuthorName:  "Vitezslav",
	// 	ImageURL:    "https://www.something.some/assets/logo-image.png",
	// 	PageURL:     "https://www.something.some",
	// }

	// content := template.HTML("<h1>Blog Post</h1><p>This is a blog post</p>")
	// blog := &models.Blog{
	// 	Slug:    fmt.Sprintf("blog-post-%d", 10000),
	// 	SEO:     seo,
	// 	Content: content,
	// }

	// new_blogs = append(new_blogs, blog)

	return &BlogHandler{
		blogs: new_blogs,
	}
}

func (bh *BlogHandler) GetAllBlogPosts(c *gin.Context) {
	page := c.Query("page")

	if page != "" {
		page_int, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid page"})
			return
		}

		start := (page_int - 1) * PER_PAGE
		end := page_int * PER_PAGE

		if start < 0 {
			start = 0
			end = 0
		}

		var show []*models.Blog
		if start < len(bh.blogs) {
			if end > len(bh.blogs) {
				show = bh.blogs[start:]
			} else {
				show = bh.blogs[start:end]
			}
		}

		c.HTML(http.StatusOK, "", view.BlogListPage(page_int+1, end < len(bh.blogs), show))

		return
	}

	c.HTML(http.StatusOK, "", view.BlogList(&models.SEO{
		Title:       "Blog",
		Description: "This is a blog",
		Keywords:    "blog",
		AuthorName:  "Vitezslav",
		ImageURL:    "https://www.something.some/assets/logo-image.png",
		PageURL:     "https://www.something.some",
	}))
}

func (bh *BlogHandler) GetBlogPost(c *gin.Context) {
	slug := c.Param("blog")

	for _, blog := range bh.blogs {
		if blog.Slug == slug {
			c.HTML(http.StatusOK, "", view.BlogDetail(blog))
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"message": "Blog post not found"})
}
