<!-- Improved compatibility of back to top link: See: https://github.com/othneildrew/Best-README-Template/pull/73 -->
<a id="readme-top"></a>
<!--
*** Thanks for checking out the Best-README-Template. If you have a suggestion
*** that would make this better, please fork the repo and create a pull request
*** or simply open an issue with the tag "enhancement".
*** Don't forget to give the project a star!
*** Thanks again! Now go create something AMAZING! :D
-->


<!-- PROJECT SHIELDS -->
<!--
*** I'm using markdown "reference style" links for readability.
*** Reference links are enclosed in brackets [ ] instead of parentheses ( ).
*** See the bottom of this document for the declaration of the reference variables
*** for contributors-url, forks-url, etc. This is an optional, concise syntax you may use.
*** https://www.markdownguide.org/basic-syntax/#reference-style-links
-->
[![LinkedIn][dhruv-linkedin-shield]][dhruv-linkedin-url]
[![Ayush's LinkedIn][ayush-linkedin-shield]][ayush-linkedin-url]


<h3 align="center">Image Beautifier</h3>

  <p align="center">
    Transform your images into masterpieces with ease and speed!
  </p>
</div>



<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
      <ul>
        <li><a href="#built-with">Built With</a></li>
      </ul>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>



<!-- ABOUT THE PROJECT -->
## About The Project

Image Beautifier makes transforming your images fast and easy! It is built in Go using Goroutines. It utilizes primitives like wait groups and design patterns like pipelines and worker pools for effective concurrency.

You can transform images in the following ways:
* Blur
* Flip upside down
* Grayscale
* Resize
* Add random cat images

<p align="right">(<a href="#readme-top">back to top</a>)</p>



### Built With

* Go

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- GETTING STARTED -->
## Getting Started

Here's how you can run Image Beautifier locally.

### Prerequisites

* Go 1.20 or newer

### Installation

1. Clone the repository
   ```sh
   git clone https://github.com/dhruvp987/ImageBeautifier.git
   ```
2. Compile the Go source into a single executable called imagebeautifer
   ```sh
   go build
   ```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- USAGE EXAMPLES -->
## Usage

```sh
imagebeautifier -i=[input image path] -o=[output image path] -c=[transformation,]
```
*NOTE:*
The input and output image must be either JPEG or PNG.

Here are the built-in transformations that you can use:
* blur
* grayscale
* resize,[multiplier]
  * [multiplier] is the multiplier to resize your image by
* upsidedown
* cats

Example:
```sh
imagebeautifier -i=myimage.png o=beautifiedimage.png -c=blur,blur,resize,3,cats,cats,upsidedown,grayscale
```

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- ROADMAP -->
## Roadmap

- [ ] Add GUI
- [ ] Create API for community transformations

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- LICENSE -->
## License

Distributed under the MIT license. See `LICENSE` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>


<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[dhruv-linkedin-shield]: https://img.shields.io/badge/dhruv's%20linkedin-blue?style=for-the-badge
[dhruv-linkedin-url]: https://linkedin.com/in/dhruvpatel789
[ayush-linkedin-shield]: https://img.shields.io/badge/ayush's%20linkedin-blue?style=for-the-badge
[ayush-linkedin-url]: https://linkedin.com/in/ayush-sharma-4320b9246
