FROM arm32v7/golang:1.22-bullseye
# FROM gocv/opencv:4.8.1-ubuntu-22.04
# FROM gocv/opencv:4.8.1
# FROM arm32v7/golang:1.22-alpine

# INSTALL BUILD DEPENDENCIES
# RUN apk update && apk add make gcc g++ pkgconfig opencv-dev g++-arm-linux-eabi gcc-arm-linux-eabi binutils-arm-linux-eabi
RUN apt-get update --allow-unauthenticated
RUN apt-get install -y sudo build-essential cmake git pkg-config libgtk-3-dev libavcodec-dev libavformat-dev libswscale-dev libv4l-dev libxvidcore-dev libx264-dev libjpeg-dev libpng-dev libtiff-dev gfortran openexr libatlas-base-dev python3-dev python3-numpy libtbb2 libtbb-dev libdc1394-22-dev
# gcc-arm-linux-gnueabihf g++-arm-linux-gnueabihf
# RUN apt-get update && apt-get install -y build-essential pkg-config gcc-10-arm-linux-gnueabihf g++-10-arm-linux-gnueabihf
# libopencv-dev 

# INSTALL OPENCV FROM DEB PACKAGE
# RUN wget https://github.com/prepkg/opencv-raspberrypi/releases/download/4.8.0/opencv.deb
# RUN dpkg -i opencv.deb

# INSTALL GO
# RUN wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
# RUN rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
# RUN export PATH=$PATH:/usr/local/go/bin

# INSTALL GOCV
RUN mkdir /gocv-download
WORKDIR /gocv-download
RUN wget -O gocv.tar.gz https://github.com/hybridgroup/gocv/archive/refs/tags/v0.35.0.tar.gz
RUN tar -xvf gocv.tar.gz
RUN mv gocv-0.35.0 gocv
WORKDIR /gocv-download/gocv
RUN make install BUILD_SHARED_LIBS=OFF

# INSTALL OPENCV
# RUN mkdir /opencv-download
# WORKDIR /opencv-download
# RUN wget -O opencv-4.9.0.tar.gz https://github.com/opencv/opencv/archive/refs/tags/4.9.0.tar.gz
# RUN tar -xvf opencv-4.9.0.tar.gz
# RUN mv opencv-4.9.0 opencv
# RUN wget -O opencv_contrib-4.9.0.tar.gz https://github.com/opencv/opencv_contrib/archive/refs/tags/4.9.0.tar.gz
# RUN tar -xvf opencv_contrib-4.9.0.tar.gz
# RUN mv opencv_contrib-4.9.0 opencv_contrib
# WORKDIR opencv
# RUN mkdir build
# WORKDIR build
# RUN cmake -D CMAKE_BUILD_TYPE=RELEASE -D CMAKE_INSTALL_PREFIX=/usr/local -D INSTALL_C_EXAMPLES=ON -D INSTALL_PYTHON_EXAMPLES=ON -D OPENCV_GENERATE_PKGCONFIG=ON -D OPENCV_EXTRA_MODULES_PATH=/opencv-download/opencv_contrib/modules -D BUILD_EXAMPLES=ON -D BUILD_SHARED_LIBS=OFF ..
# RUN make -j$(nproc)
# RUN make install
# RUN ldconfig

# INSTALL MUSL CROSS COMPILER
# RUN wget -P ~ https://musl.cc/armv7l-linux-musleabihf-cross.tgz
# RUN tar -xvf ~/armv7l-linux-musleabihf-cross.tgz -C /usr/local
# RUN ln -f -s /usr/local/armv7l-linux-musleabihf-cross/libc.so /usr/local/armv7l-linux-musleabihf-cross/armv7l-linux-musleabihf/lib/ld-musl-armhf.so.1
# RUN export PATH=$PATH:/usr/local/armv7l-linux-musleabihf-cross/bin

ENV CGO_ENABLED=1
ENV GOARCH=arm
ENV GOARM=7
ENV GOOS=linux
# ENV CGO_CPPFLAGS="-I/usr/local/include -I/usr/local/include/opencv4"
# ENV CGO_LDFLAGS="-I/usr/local/include/opencv4 -L/usr/local/lib -L/usr/local/lib/opencv4/3rdparty -lopencv_gapi -lopencv_stitching -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_dnn_objdetect -lopencv_dpm -lopencv_face -lopencv_fuzzy -lopencv_hfs -lopencv_img_hash -lopencv_line_descriptor -lopencv_quality -lopencv_reg -lopencv_rgbd -lopencv_saliency -lopencv_stereo -lopencv_structured_light -lopencv_phase_unwrapping -lopencv_superres -lopencv_optflow -lopencv_surface_matching -lopencv_tracking -lopencv_datasets -lopencv_text -lopencv_highgui -lopencv_dnn -lopencv_plot -lopencv_videostab -lopencv_video -lopencv_videoio -lopencv_xfeatures2d -lopencv_shape -lopencv_ml -lopencv_ximgproc -lopencv_xobjdetect -lopencv_objdetect -lopencv_calib3d -lopencv_imgcodecs -lopencv_features2d -lopencv_flann -lopencv_xphoto -lopencv_photo -lopencv_imgproc -lopencv_core -littnotify -llibprotobuf -lIlmImf -lquirc -lippiw -lippicv -lade -lz -ljpeg -ldl -lm -lpthread -lrt -lopenblas -lquadmath -linference_engine_s -lpugixml -lngraph -lgfortran -ltbb"
# ENV CGO_LDFLAGS="-L/usr/local/lib -lopencv_core -lopencv_imgproc -lopencv_highgui -lopencv_imgcodecs -lopencv_videoio -lopencv_video -lopencv_objdetect -lopencv_features2d -lopencv_calib3d -lopencv_flann -lopencv_ml -lopencv_dnn -lopencv_xfeatures2d -lopencv_ximgproc -lopencv_xobjdetect -lopencv_xphoto -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_datasets -lopencv_dpm -lopencv_face -lopencv_freetype -lopencv_fuzzy -lopencv_hdf -lopencv_hfs -lopencv_img_hash -lopencv_line_descriptor -lopencv_optflow -lopencv_phase_unwrapping -lopencv_plot -lopencv_reg -lopencv_rgbd -lopencv_saliency -lopencv_shape -lopencv_stereo -lopencv_structured_light -lopencv_superres -lopencv_surface_matching -lopencv_text -lopencv_tracking -lopencv_videostab -lopencv_xfeatures2d -lopencv_ximgproc -lopencv_xobjdetect -lopencv_xphoto -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_datasets -lopencv_dpm -lopencv_face -lopencv_freetype -lopencv_fuzzy -lopencv_hdf -lopencv_hfs -lopencv_img_hash -lopencv_line_descriptor -lopencv_optflow -lopencv_phase_unwrapping -lopencv_plot -lopencv_reg -lopencv_rgbd -lopencv_saliency -lopencv_shape -lopencv_stereo -lopencv_structured_light -lopencv_superres -lopencv_surface_matching -lopencv_text -lopencv_tracking -lopencv_videostab -lopencv_xfeatures2d -lopencv_ximgproc -lopencv_xobjdetect -lopencv_xphoto -lopencv_aruco -lopencv_bgsegm -lopencv_bioinspired -lopencv_ccalib -lopencv_datasets -lopencv_dpm -lopencv_face -lopencv_freetype -lopencv_fuzzy -lopencv_hdf -lopencv_hfs -lopencv_img_hash -lopencv_line_descriptor -lopencv_optflow -lopencv_phase_unwrapping -lopencv_plot -lopencv"
# -lopencv_gapi
# -lopencv_stitching
# -lopencv_aruco
# -lopencv_bgsegm
# -lopencv_bioinspired
# -lopencv_ccalib
# -lopencv_dnn_objdetect
# -lopencv_dnn_superres
# -lopencv_dpm
# -lopencv_face
# -lopencv_freetype
# -lopencv_fuzzy
# -lopencv_hfs
# -lopencv_img_hash
# -lopencv_intensity_transform
# -lopencv_line_descriptor
# -lopencv_mcc
# -lopencv_quality
# -lopencv_rapid
# -lopencv_reg
# -lopencv_rgbd
# -lopencv_saliency
# -lopencv_stereo
# -lopencv_structured_light
# -lopencv_phase_unwrapping
# -lopencv_superres
# -lopencv_optflow
# -lopencv_surface_matching
# -lopencv_tracking
# -lopencv_highgui
# -lopencv_datasets
# -lopencv_text
# -lopencv_plot
# -lopencv_videostab
# -lopencv_videoio
# -lopencv_wechat_qrcode
# -lopencv_xfeatures2d
# -lopencv_shape
# -lopencv_ml
# -lopencv_ximgproc
# -lopencv_video
# -lopencv_xobjdetect
# -lopencv_objdetect
# -lopencv_calib3d
# -lopencv_imgcodecs
# -lopencv_features2d
# -lopencv_dnn
# -lopencv_flann
# -lopencv_xphoto
# -lopencv_photo
# -lopencv_imgproc
# -lopencv_core
# ENV CC=/usr/local/armv7l-linux-musleabihf-cross/bin/armv7l-linux-musleabihf-gcc
# ENV CXX=/usr/local/armv7l-linux-musleabihf-cross/bin/armv7l-linux-musleabihf-g++

RUN mkdir /app
COPY . /app
WORKDIR /app

RUN go mod download

# CMD ["go", "build", "-v", "-trimpath", "-ldflags", "-extldflags -static"]
# RUN go build -v -trimpath -ldflags="-extldflags -static"
RUN go build --ldflags '-extldflags "-static"'

CMD ["echo", "Build complete!"]




