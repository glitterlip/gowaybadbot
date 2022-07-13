//定义Points的构造函数
(function ($, window, document, undefined) {
    var Points = function (ele, opt) {
        this.$element = ele
        this.defaults = {
            mode: 'fixed',	//弹出式pop，固定fixed
            defaultNum: 4,	//默认的文字数量
            checkNum: 2,	//校对的文字数量
            vSpace: 5,	//间隔
            imgSize: {
                width: '320px',
                height: '240px',
            },
            barSize: {
                width: '320px',
                height: '40px',
            },
            hints: [],
            ready: function () {
            },
            success: function () {
            },
            error: function () {
            }
        };
        this.options = $.extend({}, this.defaults, opt)
    };

//定义Points的方法
    Points.prototype = {
        init: function () {
            const _this = this;
            //加载页面
            _this.loadDom();
            _this.refresh();
            _this.options.ready();
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
            //点击事件比对
            _this.$element.find('.verify-img-panel img').on('click', function (e) {
                _this.checkPosArr.push(_this.getMousePos(this, e));
                if (_this.num === _this.options.checkNum - 1) {
                    _this.createPoint(_this.getMousePos(this, e));
                    _this.check(_this);
                } else {
                    _this.createPoint(_this.getMousePos(this, e));
                    let next = _this.hints[_this.num]
                    _this.htmlDoms.msg.text(`请点击[${next}]`)
                }
            });
            //刷新
            _this.$element.find('.verify-refresh').on('click', function () {
                _this.refresh(true);
            });
        },
        //加载页面
        loadDom: function () {
            // this.fontPos = [];	//选中的坐标信息
            this.checkPosArr = [];	//用户点击的坐标
            this.num = 0;	//点击的记数
            this.setSize = resetSize(this);	//重新设置宽度高度
            this.htmlDoms = {
                out_panel: this.$element.find('.verify-img-out'),
                img_panel: this.$element.find('.verify-img-panel'),
                bar_area: this.$element.find('.verify-bar-area'),
                msg: this.$element.find('.verify-msg'),
            };
            this.initHint = this.htmlDoms.msg.text()
            this.hints = this.htmlDoms.msg.text().match(/\[(.*)]/)[1].split(',')
            this.$element.css('position', 'relative');
            if (this.options.mode === MODE_POP) {
                this.htmlDoms.out_panel.css({'display': 'none', 'position': 'absolute', 'bottom': '42px'});
            } else {
                this.htmlDoms.out_panel.css({'position': 'relative'});
            }

            this.htmlDoms.out_panel.css('height', parseInt(this.setSize.img_height) + this.options.vSpace + 'px');
            this.htmlDoms.img_panel.css({
                'width': this.setSize.img_width,
                'height': this.setSize.img_height,
                'background-size': this.setSize.img_width + ' ' + this.setSize.img_height,
                'margin-bottom': this.options.vSpace + 'px'
            });
            this.htmlDoms.bar_area.css({
                'width': this.options.barSize.width,
                'height': this.setSize.bar_height,
                'line-height': this.setSize.bar_height
            });

        },
        //获取坐标
        getMousePos: function (obj, event) {
            var e = event || window.event;
            var scrollX = document.documentElement.scrollLeft || document.body.scrollLeft;
            var scrollY = document.documentElement.scrollTop || document.body.scrollTop;
            var x = e.clientX - ($(obj).offset().left - $(window).scrollLeft());
            var y = e.clientY - ($(obj).offset().top - $(window).scrollTop());

            return {'x': x, 'y': y};
        },

        //创建坐标点
        createPoint: function (pos) {
            this.htmlDoms.img_panel.append('<div class="point-area" style="background-color:#1abd6c;color:#fff;z-index:9999;width:30px;height:30px;text-align:center;line-height:30px;border-radius: 50%;position:absolute;top:' + parseInt(pos.y - 10) + 'px;left:' + parseInt(pos.x - 10) + 'px;">' + this.num + '</div>');
            this.num++;
        },

        //弹出式
        showImg: function () {
            this.htmlDoms.out_panel.css({'display': 'block'});
        },

        //固定式
        hideImg: function () {
            this.htmlDoms.out_panel.css({'display': 'none'});
        },

        //刷新
        refresh: function (reload = false) {
            var _this = this;
            this.$element.find('.point-area').remove();
            this.fontPos = [];
            this.checkPosArr = [];
            this.num = 0;
            _this.$element.find('.verify-bar-area').css({'color': '#000', 'border-color': '#ddd'});
            _this.$element.find('.verify-refresh').show();
            if (reload) {
                $.ajax({
                    "url": "/refresh",
                    "async": false,
                    "method": "GET",
                    "data": {"old-challenge": $("#challenge-id").val()}
                }).then((r) => {
                    $('#challenge-id').val(r.meta.id)
                    $('.verify-img-canvas').attr("src", r.meta.rule.rule_image)
                    $('.verify-msg').text(r.meta.rule.hint)
                    _this.initHint = r.meta.rule.hint
                    _this.hints = this.htmlDoms.msg.text().match(/\[(.*)]/)[1].split(',')
                    _this.$element.find('.verify-msg').text(_this.initHint);
                })
            }
        },
        check: (_this) => {
            let positions = _this.checkPosArr
            let answer = `${positions[0].x},${positions[0].y},${positions[1].x},${positions[1].y}`
            $.ajax({
                "url": "/check",
                "async": false,
                "method": "POST",
                "dataType": "json",
                "data": {"challenge": $("#challenge-id").val(), answer}
            }).then((r) => {
                switch (r.code) {
                    case 0:
                        _this.success(_this)
                        break;
                    case 2:
                        _this.refresh(true)
                        break;
                    case 3:
                        _this.fail(_this)
                }
            })
        },
        success: (_this) => {
            _this.$element.find('.verify-bar-area').css({'color': '#4cae4c', 'border-color': '#5cb85c'});
            _this.$element.find('.verify-msg').text('验证成功');
            _this.$element.find('.verify-refresh').hide();
            _this.$element.find('.verify-img-panel').unbind('click');
            _this.options.success(_this);
        },
        fail: (_this) => {
            _this.options.error(_this);
            _this.$element.find('.verify-bar-area').css({'color': '#d9534f', 'border-color': '#d9534f'});
            _this.$element.find('.verify-msg').text("验证失败");

            setTimeout(function () {
                _this.$element.find('.verify-msg').text(_this.initHint);
                _this.$element.find('.verify-bar-area').css({'color': '#000', 'border-color': '#ddd'});
                _this.refresh();
            }, 1000);
        }

    };


//在插件中使用clickVerify对象
    $.fn.pointsVerify = function (options, callbacks) {
        var points = new Points(this, options);
        points.init();
    };
})(jQuery, window, document);