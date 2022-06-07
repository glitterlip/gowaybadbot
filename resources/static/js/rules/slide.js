const CODE_TYPE_BASE = 1; //普通验证码

const SLIDE_TYPE_BASE = 1;
const SLIDE_TYPE_SEAL = 2;
const SLIDE_MODE_FIXED = "FIXED";

(function ($, window, document, undefined) {
    var Slide = function (ele, opt) {
        this.$element = ele;//recaptcha html element
        this.defaultOptions = {
            type: SLIDE_TYPE_BASE,
            mode: SLIDE_MODE_FIXED,	//弹出式pop，固定fixed
            vOffset: 5, //accuracy px
            vSpace: 5,
            explain: '向右滑动完成验证',
            imgSize: {  //slide background image size
                width: '320px',
                height: '240px',
            },
            blockSize: { //slide seal image size
                width: '50px',
                height: '50px',
            },
            circleRadius: '10px',
            barSize: { //slide bar size
                width: '320px',
                height: '40px',
            },
            ready: function () {
            },
            success: function () {
            },
            error: function () {
            }

        }
        this.options = $.extend({}, this.defaultOptions, opt)
    };
    //定义Slide的方法
    Slide.prototype = {
        init: function () {
            const _this = this;

            //加载页面
            this.loadDom();
            _this.refresh();
            this.options.ready();

            this.$element[0].onselectstart = document.body.ondrag = function () {
                return false;
            };

            if (this.options.mode === MODE_POP) {
                this.$element.on('mouseover', function (e) {
                    _this.showImg();
                });

                this.$element.on('mouseout', function (e) {
                    _this.hideImg();
                });


                this.htmlDoms.out_panel.on('mouseover', function (e) {
                    _this.showImg();
                });

                this.htmlDoms.out_panel.on('mouseout', function (e) {
                    _this.hideImg();
                });
            }

            //按下
            this.htmlDoms.move_block.on('touchstart', function (e) {
                _this.start(e);
            });

            this.htmlDoms.move_block.on('mousedown', function (e) {
                _this.start(e);
            });

            //拖动
            window.addEventListener("touchmove", function (e) {
                _this.move(e);
            });


            window.addEventListener("mousemove", function (e) {

                _this.move(e);
            });

            //鼠标松开
            window.addEventListener("touchend", function () {
                _this.end();
            });
            window.addEventListener("mouseup", function () {
                _this.end();
            });

            //刷新
            _this.$element.find('.verify-refresh').on('click', function () {
                _this.refresh(true);
            });
        },

        //初始化加载
        loadDom: function () {
            this.status = false;	//verify status
            this.isFinished = false;		//是否完成验证
            this.sizeConfig = this.resetSize(this);	//重新设置宽度高度
            this.plusWidth = 0;//seal width
            this.plusHeight = 0;//seal height
            this.x = 0;
            this.y = 0;
            this.lengthPercent = (parseInt(this.sizeConfig.img_width) - parseInt(this.sizeConfig.block_width) - parseInt(this.sizeConfig.circle_radius) - parseInt(this.sizeConfig.circle_radius) * 0.8) / (parseInt(this.sizeConfig.img_width) - parseInt(this.sizeConfig.bar_height));


            this.htmlDoms = {
                stamp_image: this.$element.find('.verify-sub-block'),
                out_panel: this.$element.find('.verify-img-out'),
                img_panel: this.$element.find('.verify-img-panel'),
                background_image: this.$element.find('.verify-img-background'),
                bar_area: this.$element.find('.verify-bar-area'),
                move_block: this.$element.find('.verify-move-block'),
                left_bar: this.$element.find('.verify-left-bar'),
                msg: this.$element.find('.verify-msg'),
                icon: this.$element.find('.verify-icon'),
                refresh: this.$element.find('.verify-refresh')
            };
            this.$element.css('position', 'relative');
            if (this.options.mode === MODE_POP) {
                this.htmlDoms.out_panel.css({'display': 'none', 'position': 'absolute', 'bottom': '42px'});
                this.htmlDoms.stamp_image.css({'display': 'none'});
            } else {
                this.htmlDoms.out_panel.css({'position': 'relative'});
            }

            this.htmlDoms.out_panel.css('height', parseInt(this.sizeConfig.img_height) + this.options.vSpace + 'px');
            this.htmlDoms.img_panel.css({'width': this.sizeConfig.img_width, 'height': this.sizeConfig.img_height});
            this.htmlDoms.bar_area.css({
                'width': this.sizeConfig.bar_width,
                'height': this.sizeConfig.bar_height,
                'line-height': this.sizeConfig.bar_height
            });
            this.htmlDoms.move_block.css({'width': "40px", 'height': this.sizeConfig.bar_height});
            this.htmlDoms.left_bar.css({'width': "40px", 'height': this.sizeConfig.bar_height});
            this.randSet();
        },
        //鼠标按下
        start: function (e) {
            if (this.isFinished === false) {
                this.htmlDoms.msg.text('');
                this.htmlDoms.move_block.css('background-color', '#337ab7');
                this.htmlDoms.left_bar.css('border-color', '#337AB7');
                this.htmlDoms.icon.css('color', '#fff');
                e.stopPropagation();
                this.status = true;
            }
        },

        //鼠标移动
        move: function (e) {
            if (this.status && this.isFinished === false) {
                if (this.options.mode === MODE_POP) {
                    this.showImg();
                }

                if (!e.touches) {    //兼容移动端
                    var x = e.clientX;
                } else {     //兼容PC端
                    var x = e.touches[0].pageX;
                }
                const bar_area_left = Slide.prototype.getLeft(this.htmlDoms.bar_area[0]);
                let move_block_left = x - bar_area_left; //小方块相对于父元素的left值


                if (move_block_left >= (this.htmlDoms.bar_area[0].offsetWidth - parseInt(this.sizeConfig.bar_height) + parseInt(parseInt(this.sizeConfig.block_width) / 2) - 2)) {
                    move_block_left = (this.htmlDoms.bar_area[0].offsetWidth - parseInt(this.sizeConfig.bar_height) + parseInt(parseInt(this.sizeConfig.block_width) / 2) - 2);
                }

                if (move_block_left <= parseInt(parseInt(this.sizeConfig.block_width) / 2)) {
                    move_block_left = parseInt(parseInt(this.sizeConfig.block_width) / 2);
                }


                //拖动后小方块的left值
                this.htmlDoms.move_block.css('left', move_block_left - parseInt(parseInt(this.sizeConfig.block_width) / 2) + "px");
                this.htmlDoms.left_bar.css('width', move_block_left - parseInt(parseInt(this.sizeConfig.block_width) / 2) + "px");
                let t = (move_block_left - parseInt(parseInt(this.sizeConfig.block_width) / 2)) + "px"
                this.htmlDoms.stamp_image.css('left', t);

            }
        },

        //鼠标松开
        end: function () {
            var _this = this;
            //判断是否重合
            if (this.status && this.isFinished === false) {
                if (this.options.type === SLIDE_TYPE_SEAL) {		//图片滑动
                    let answer = parseInt(this.htmlDoms.stamp_image.css('left'))
                    $.ajax({
                        "url": "/check",
                        "async": false,
                        "method": "POST",
                        "dataType": "json",
                        "data": {"challenge": $("#challenge-id").val(), answer}
                    }).then((r) => {
                        switch (r.code) {
                            case 0:
                                this.success(_this)
                                break;
                            case 2:
                                this.refresh(true)
                                break;
                            case 3:
                                this.fail(_this)
                        }
                    })
                }
                this.status = false
            }
        },
        success: (_this) => {
            _this.htmlDoms.move_block.css('background-color', '#5cb85c');
            _this.htmlDoms.left_bar.css({'border-color': '#5cb85c', 'background-color': '#fff'});
            _this.htmlDoms.icon.css('color', '#fff');
            _this.htmlDoms.icon.removeClass('icon-right');
            _this.htmlDoms.icon.addClass('icon-check');
            _this.htmlDoms.refresh.hide();
            _this.isFinished = true;
            _this.options.success(this);
        },
        fail: (_this) => {
            _this.htmlDoms.move_block.css('background-color', '#d9534f');
            _this.htmlDoms.left_bar.css('border-color', '#d9534f');
            _this.htmlDoms.icon.css('color', '#fff');
            _this.htmlDoms.icon.removeClass('icon-right');
            _this.htmlDoms.icon.addClass('icon-close');
            setTimeout(function () {
                _this.refresh();
            }, 400);

            _this.options.error(this);
        },
        //弹出式
        showImg: function () {
            this.htmlDoms.out_panel.css({'display': 'block'});
            this.htmlDoms.stamp_image.css({'display': 'block'});
        },
        //固定式
        hideImg: function () {
            this.htmlDoms.out_panel.css({'display': 'none'});
            this.htmlDoms.stamp_image.css({'display': 'none'});
        },
        //calculate element size
        resetSize: function (obj) {
            let img_width, img_height, bar_width, bar_height, block_width, block_height, circle_radius;	//图片的宽度、高度，移动条的宽度、高度
            const parentWidth = obj.$element.parent().width() || $(window).width();
            const parentHeight = obj.$element.parent().height() || $(window).height();
            //eg: if container size is 100px, and image size is 50%,so the image real size is 50px
            if (obj.options.imgSize.width.indexOf('%') !== -1) {
                img_width = parseInt(obj.options.imgSize.width) / 100 * parentWidth + 'px';
            } else {
                img_width = obj.options.imgSize.width;
            }
            if (obj.options.imgSize.height.indexOf('%') !== -1) {
                img_height = parseInt(obj.options.imgSize.height) / 100 * parentHeight + 'px';
            } else {
                img_height = obj.options.imgSize.height;
            }
            if (obj.options.barSize.width.indexOf('%') !== -1) {
                bar_width = parseInt(obj.options.barSize.width) / 100 * parentWidth + 'px';
            } else {
                bar_width = obj.options.barSize.width;
            }
            if (obj.options.barSize.height.indexOf('%') !== -1) {
                bar_height = parseInt(obj.options.barSize.height) / 100 * parentHeight + 'px';
            } else {
                bar_height = obj.options.barSize.height;
            }
            if (obj.options.blockSize) {
                if (obj.options.blockSize.width.indexOf('%') !== -1) {
                    block_width = parseInt(obj.options.blockSize.width) / 100 * parentWidth + 'px';
                } else {
                    block_width = obj.options.blockSize.width;
                }
                if (obj.options.blockSize.height.indexOf('%') !== -1) {
                    block_height = parseInt(obj.options.blockSize.height) / 100 * parentHeight + 'px';
                } else {
                    block_height = obj.options.blockSize.height;
                }
            }
            return {
                img_width: img_width,
                img_height: img_height,
                bar_width: bar_width,
                bar_height: bar_height,
                block_width: block_width,
                block_height: block_height,
            };
        },

        //随机出生点位
        randSet: function () {
            // const rand1 = Math.floor(Math.random() * 9 + 1);
            // const rand2 = Math.floor(Math.random() * 9 + 1);
            // const top = rand1 * parseInt(this.sizeConfig.img_height) / 15 + parseInt(this.sizeConfig.img_height) * 0.1;
            // const left = rand2 * parseInt(this.sizeConfig.img_width) / 15 + parseInt(this.sizeConfig.img_width) * 0.1;
            //
            // this.x = left;
            // this.y = top;
            if (this.options.mode === MODE_POP) {
                this.htmlDoms.stamp_image.css({'top': '-' + (parseInt(this.sizeConfig.img_height) + this.options.vSpace - this.y - 2) + 'px'});
            } else {
                this.htmlDoms.stamp_image.css({'top': '50px'});
            }

        },

        //刷新
        refresh: function (reload = false) {
            var _this = this;
            this.htmlDoms.refresh.show();
            this.$element.find('.verify-msg:eq(1)').text('');
            this.$element.find('.verify-msg:eq(1)').css('color', '#000');
            this.htmlDoms.move_block.animate({'left': '0px'}, 'fast');
            this.htmlDoms.left_bar.animate({'width': parseInt(this.sizeConfig.bar_height)}, 'fast');
            this.htmlDoms.left_bar.css({'border-color': '#ddd'});
            this.htmlDoms.move_block.css('background-color', '#fff');
            this.htmlDoms.icon.css('color', '#000');
            this.htmlDoms.icon.removeClass('icon-close');
            this.htmlDoms.icon.addClass('icon-right');
            this.$element.find('.verify-msg:eq(0)').text(this.options.explain);
            this.randSet();
            this.isFinished = false;
            this.htmlDoms.stamp_image.css('left', "0px");
            if (reload) {
                $.ajax({
                    "url": "/refresh",
                    "async": false,
                    "method": "GET",
                    "data": {"old-challenge": $("#challenge-id").val()}
                }).then((r) => {
                    $('#challenge-id').val(r.meta.Id)
                    $('.verify-img-canvas').attr("src", r.meta.Rule.BackgroundImage)
                    $('.verify-sub-block').attr("src", r.meta.Rule.SealImage)
                })
            }
        },

        //获取left值
        getLeft: function (node) {
            var left = $(node).offset().left;
            return left;
        },
    };

    $.fn.slideVerify = function (options, callbacks) {
        const slide = new Slide(this, options);
        slide.init();
    };

})(jQuery, window, document);