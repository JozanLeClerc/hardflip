# ========================
# =====    ===============
# ======  ================
# ======  ================
# ======  ====   ====   ==
# ======  ===     ==  =  =
# ======  ===  =  ==     =
# =  ===  ===  =  ==  ====
# =  ===  ===  =  ==  =  =
# ==     =====   ====   ==
# ========================
#
# hardflip: Makefile
# Tue Jan 23 11:16:43 2024
# Joe
#
# GNU Makefile

TARGET			:= hf
SHELL			:= /bin/sh
SRC_DIR			:= ./src/
SRC_NAME		:= *.go
SRC				 = $(addprefix ${SRC_DIR}, ${SRC_NAME})
.DEFAULT_GOAL	:= ${TARGET}

run: ${SRC}
	go run ${SRC_DIR}

${TARGET}: ${SRC}
	go build -o ${TARGET} ${SRC_DIR}

clean:
	go clean
	rm -f ${TARGET}

.PHONY: hf run clean
