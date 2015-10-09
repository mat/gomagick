# gomagick

Image resizing server written in Go based on ImageMagick, proof of concept.

See in action at <http://gomagick.herokuapp.com>

Powered by ImageMagick bindings from <https://github.com/quirkey/magick>

## Usage

Example: <http://gomagick.herokuapp.com/img?size=100x100&url=https://github.com/apple-touch-icon.png>


* Endpoint: <http://gomagick.herokuapp.com/img>
* Params:
	* `url`: http/https URL
	* `size`: 64x64, 90%



## Performance

Every request returns the `X-Image-Timings` header containing performance metrics:

	curl -I 'http://gomagick.herokuapp.com/img?size=100x100&url=https://github.com/apple-touch-icon.png'

For example:

	HTTP/1.1 200 OK
	Connection: keep-alive
	Content-Type: image/png
	X-Image-Timings: get=37.673327ms, resize=10.157948ms
	Date: Thu, 08 Oct 2015 18:41:59 GMT



## Docker

	git clone https://github.com/mat/gomagick.git
	cd gomagick
	docker build -t magickserver .
	docker run -it -p 8080:8080 --name magickserver magickserver
	open http://$(docker-machine ip default):8080

## Thanks

I started researching this during our student exchange between <https://www.xing.com> and <http://www.jimdo.com>. Thanks to both for letting me work on interesting things.

## License

The MIT License (MIT)

Copyright (c) 2015 Matthias Lüdtke, Hamburg - http://github.com/mat

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
