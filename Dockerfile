FROM ubuntu:20.04

WORKDIR /tmp

# No apt prompts during build!
RUN ln -fs /usr/share/zoneinfo/America/New_York /etc/localtime
RUN export DEBIAN_FRONTEND=noninteractive

RUN curl -fsSL https://deb.nodesource.com/setup_14.x | bash - \
    && apt-get update -y && apt-get upgrade -y \
    && apt-get install -y wget git build-essential nodejs npm \
    # puppeteer dependencies
    ca-certificates fonts-liberation libappindicator3-1 libasound2 libatk-bridge2.0-0 libatk1.0-0 \
    libc6 libcairo2 libcups2 libdbus-1-3 libexpat1 libfontconfig1 libgbm1 libgcc1 libglib2.0-0 \
    libgtk-3-0 libnspr4 libnss3 libpango-1.0-0 libpangocairo-1.0-0 libstdc++6 libx11-6 libx11-xcb1 \
    libxcb1 libxcomposite1 libxcursor1 libxdamage1 libxext6 libxfixes3 libxi6 libxrandr2 libxrender1 \
    libxss1 libxtst6 lsb-release wget xdg-utils 


# install chrome
# RUN wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb \
#     && apt update -y && apt install -y ./google-chrome-stable_current_amd64.deb

# install golang
RUN wget https://golang.org/dl/go1.16.5.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go1.16.5.linux-amd64.tar.gz \
    && echo export 'PATH=$PATH:/usr/local/go/bin' >>.profile


WORKDIR ${GOPATH}
RUN /usr/local/go/bin/go get -v github.com/ramya-rao-a/go-outline \
    && /usr/local/go/bin/go get -v github.com/uudashr/gopkgs/v2/cmd/gopkgs \
    && /usr/local/go/bin/go get -v github.com/cweill/gotests/gotests \
    && /usr/local/go/bin/go get -v github.com/fatih/gomodifytags \
    && /usr/local/go/bin/go get -v github.com/josharian/impl \
    && /usr/local/go/bin/go get -v github.com/haya14busa/goplay/cmd/goplay \
    && /usr/local/go/bin/go get -v github.com/go-delve/delve/cmd/dlv \
    && /usr/local/go/bin/go get -v honnef.co/go/tools/cmd/staticcheck \
    && /usr/local/go/bin/go get -v golang.org/x/tools/gopls

WORKDIR /go/src/github.com/thomasrea0113/gpu-price-monitor

COPY package.json .
RUN npm install -D

COPY . . 