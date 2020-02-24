package img

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"github.com/nfnt/resize"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math"
	"myTestGo/exception"
	"net/http"
	"os"
	"strings"
)

// 圆形蒙版
func circleMask(d int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, d, d))

	for x := 0; x < d; x++ {
		for y := 0; y < d; y++ {
			dis := math.Sqrt(math.Pow(float64(x-d/2), 2) + math.Pow(float64(y-d/2), 2))
			if dis > float64(d)/2 {
				img.Set(x, y, color.RGBA{255, 255, 255, 0})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 255, 255})
			}
		}
	}
	return img
}

// 将图片转为圆形
func circlePhoto(img image.Image, diameter int) (image.Image) {
	// 缩放图片
	img = resize.Resize(uint(diameter), uint(diameter), img, resize.Lanczos3)
	// 圆形蒙版
	maskImg := circleMask(diameter)
	// 圆形图片画布
	circleImg := image.NewRGBA(image.Rect(0, 0, diameter, diameter))
	// 将 img 用 maskImg 圆形蒙版画在 circleImg 上；
	// dst:背景图/画布 r:背景图的画图区域 src:要画的图 sp:对src的画图开始坐标 mask:使用的蒙版 mp:蒙版的开始坐标
	draw.DrawMask(circleImg, img.Bounds().Add(image.ZP), img, image.ZP, maskImg, image.ZP, draw.Src)
	return circleImg
}

func roundRectPhoto(src image.Image, r int) (image.Image) {
	b := src.Bounds()
	x := b.Dx()
	y := b.Dy()
	dst := image.NewRGBA(b)

	p1 := image.Point{r, r}
	p2 := image.Point{x - r, r}
	p3 := image.Point{r, y - r}
	p4 := image.Point{x - r, y - r}

	for m := 0; m < x; m++ {
		for n := 0; n < y; n++ {
			if (p1.X-m)*(p1.X-m)+(p1.Y-n)*(p1.Y-n) > r*r && m <= p1.X && n <= p1.Y {
			} else if (p2.X-m)*(p2.X-m)+(p2.Y-n)*(p2.Y-n) > r*r && m > p2.X && n <= p2.Y {
			} else if (p3.X-m)*(p3.X-m)+(p3.Y-n)*(p3.Y-n) > r*r && m <= p3.X && n > p3.Y {
			} else if (p4.X-m)*(p4.X-m)+(p4.Y-n)*(p4.Y-n) > r*r && m > p4.X && n > p4.Y {
			} else {
				dst.Set(m, n, src.At(m, n))
			}
		}
	}
	return dst
}

func ImageDraw() error {
	// 缩放比例 原始比例 300px-180px
	//zoom := 2
	// 字体 STHeiti Light.ttc   WeiRuanZhengHeiTi-2.ttc
	//fontName := "WeiRuanZhengHeiTi-2.ttc"
	//// 标题大小
	//titleSize := float64(11)
	//// 其他字体大小
	//otherFontSize := float64(9)

	// 模版底图
	tempFile, _ := os.Open("xcx.png")
	defer tempFile.Close()

	tempImg, err := png.Decode(tempFile)
	if err != nil {
		return exception.NewException(500, err.Error())
	}

	// 画布
	img := image.NewNRGBA(image.Rect(0, 0, tempImg.Bounds().Max.X, tempImg.Bounds().Max.Y))

	// 白底
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{R: 255, G: 255, B: 255, A: 255}}, image.ZP, draw.Src)

	// 将模版图画入画布
	draw.Draw(img, img.Bounds(), tempImg, image.ZP, draw.Src)

	// 获取网络图片 https://test-imgqn.smm.cn/test/usercenter/avatar/nNaGF20190923194241.jpeg
	avatarImg, err := getImg("https://test-imgqn.smm.cn/test/b/image/JaEJlzFXhZzldeWQNaya20191024055744.jpg?imageView2/1/w/150/h/150/q/100")
	if err != nil {
		return exception.NewException(500, err.Error())
	}
	// 头像转为圆形
	avatar := circlePhoto(avatarImg, 194)

	// 头像转为正方形
	//avatar := resize.Resize(39, 39, avatarImg, resize.Lanczos3)
	//头像转为圆角
	//avatar = roundRectPhoto(avatar, 4)

	// 将头像画入img
	offset := image.Pt(117, 117)
	draw.Draw(img, img.Bounds().Add(offset), avatar, image.ZP, draw.Over)

	// 画笔
	//c := freetype.NewContext()
	//
	//// 分辨率
	//c.SetDPI(100)
	//
	//// 设置背景
	//c.SetClip(img.Bounds())
	//
	//// 设置目标图像
	//c.SetDst(img)
	//c.SetHinting(font.HintingFull)
	//
	//// 加载字体
	//fontFam, err := getFontFamily(fontName)
	//if err != nil {
	//	return exception.NewException(500, err.Error())
	//}
	//// 设置字体
	//c.SetFont(fontFam)
	//
	//// 设置字体颜色
	//c.SetSrc(image.NewUniform(color.RGBA{R: 0, G: 33, B: 48, A: 255}))
	//
	//// 字体大小
	//c.SetFontSize(titleSize)
	//
	//// 写字
	//_, err = c.DrawString("张新杰", freetype.Pt(80, 40))
	//if err != nil {
	//	return exception.NewException(500, err.Error())
	//}
	//
	//c.SetSrc(image.NewUniform(color.RGBA{R: 118, G: 141, B: 154, A: 255}))
	//c.SetFontSize(otherFontSize)
	//_, err = c.DrawString("工程师", freetype.Pt(160, 48))
	//if err != nil {
	//	return exception.NewException(500, err.Error())
	//}
	//_, err = c.DrawString("上海有色网", freetype.Pt(80, 65))
	//if err != nil {
	//	return exception.NewException(500, err.Error())
	//}
	//_, err = c.DrawString("上海市浦东新区陆家嘴软件园9号楼", freetype.Pt(40, 117))
	//if err != nil {
	//	return exception.NewException(500, err.Error())
	//}
	//
	////c.SetSrc(image.NewUniform(color.RGBA{R: 0, G: 125, B: 177, A: 255}))
	//_, err = c.DrawString("13127881019", freetype.Pt(40, 92))
	//if err != nil {
	//	return exception.NewException(500, err.Error())
	//}

	// 图片转字节数组
	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return exception.NewException(500, err.Error())
	}
	imgData := buf.Bytes()

	// 转图片文件
	filename := "dst.png"
	file, _ := os.Create(filename)
	defer file.Close()
	file.Write(imgData)

	return nil
}

// 获取字符集，仅调用一次
func getFontFamily(fontName string) (*truetype.Font, error) {
	// 这里需要读取中文字体，否则中文文字会变成方格
	fontBytes, err := ioutil.ReadFile(fontName)
	if err != nil {
		fmt.Println("read file error:", err)
		return &truetype.Font{}, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println("parse font error:", err)
		return &truetype.Font{}, err
	}

	return f, err
}

func getImg(url string) (img image.Image, err error) {
	path := strings.Split(url, "/")
	var name string
	if len(path) > 1 {
		name = path[len(path)-1]
	}
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if strings.LastIndex(name, "png") > -1 {
		return png.Decode(resp.Body)
	} else {
		return jpeg.Decode(resp.Body)
	}
}
