const CODE_TYPE_BASE = 1; //普通验证码


(function ($, window, document, undefined) {
    //定义Code的构造函数
    var Code = function (ele, opt) {
        this.$element = ele;
        this.defaults = {
            type: CODE_TYPE_BASE, //普通验证码类型
            figure: 100,	//位数，仅在验证码类型为计算型时生效 默认100以内加减乘
            arith: 0,	//算法，支持加减乘，0为随机，仅在验证码类型为计算型时生效
            width: '200px',
            height: '100px',
            fontSize: '30px',
            codeLength: 6,//verifycode length only works when verify code type is CODE_TYPE_BASE
            btnId: 'check-btn',
            ready: function () {
            },
            success: function () {
            },
            error: function () {
            }
        }
        //merge options
        this.options = $.extend({}, this.defaults, opt)
    };
    //定义Code的方法
    Code.prototype = {
        init: function () {
            const _this = this;
            this.loadDom();
            this.options.ready();

            this.$element[0].onselectstart = document.body.ondrag = function () {
                return false;
            };

            //点击验证码
            this.$element.find('.verify-code, .verify-change-code').on('click', function () {

                _this.setCode(_this);

            });

            //确定的点击事件
            this.htmlDoms.code_btn.on('click', function () {
                _this.checkCode(_this);
            })

        },

        //加载页面
        loadDom: function () {
            this.isFinished = false;
            this.htmlDoms = {
                code_btn: $('#' + this.options.btnId),
                code: this.$element.find('.verify-code'),
                code_area: this.$element.find('.verify-code-area'),
                code_input: this.$element.find('.varify-input-code'),
            };

            this.htmlDoms.code.css({
                'width': this.options.width,
                'height': this.options.height,
                'line-height': this.options.height,
                'font-size': this.options.fontSize
            });
            this.htmlDoms.code_area.css({'width': this.options.width});
        },


        //设置验证码
        setCode: function () {
            if (this.isFinished === false) {
                $.ajax({
                    "url": "/refresh",
                    "async": false,
                    "method": "GET",
                    "data": {"old-challenge": $("#challenge-id").val()}
                }).then((r) => {
                    $('#challenge-id').val(r.meta.Id)
                    $('#rule-image').attr("src", r.meta.Rule.RuleImage)
                })
            }
        },

        //比对验证码
        checkCode: function (_this) {
            let answer;
            answer = this.htmlDoms.code_input.val().toUpperCase();
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
                    case 3:
                        this.fail(_this)
                }
            })

        },
        success: (_this) => {
            _this.isFinished = true;
            _this.options.success(_this);
        },
        fail: (_this) => {
            _this.options.error(_this);
            $('#hint').text('验证码错误，请输入新验证码')
            _this.htmlDoms.code_input.val()
            _this.setCode();
        },
    };

    $.fn.codeVerify = function (options, callbacks) {
        const code = new Code(this, options);
        code.init();
    };


})(jQuery, window, document);