# Executable name
TARGET = bin/inblog
SIDE_EFFECTS = .cache content public bin

# Go parameters
GOCMD = go
GOBUILD = $(GOCMD) build
GOCLEAN = $(GOCMD) clean

# Platforms
PLATFORMS = linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64 windows/arm64

all: $(TARGET)

$(TARGET):
	$(GOBUILD) -o $(TARGET) .

all-platforms:
	$(foreach platform,$(PLATFORMS),\
		$(eval OS := $(word 1,$(subst /, ,$(platform))))\
		$(eval ARCH := $(word 2,$(subst /, ,$(platform))))\
		GOOS=$(OS) GOARCH=$(ARCH) $(GOBUILD) -ldflags '-w' -o $(TARGET)-$(OS)-$(ARCH) .;)

test:
	$(GOCMD) test .

clean:
	$(GOCLEAN)
	rm -rf $(SIDE_EFFECTS) $(TARGET)

.PHONY: all clean all-platforms
