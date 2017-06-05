var Panel = function () {

    var createHandler = function () {
        console.log("init called...")
        $('#create-btn').on('click', function (e) {
            e.preventDefault();
            console.log("create-btn clicked...");

            //1. Hide create-btn.
            $('#form-title').css('display', 'none');
            $('#create-btn').css('display', 'none');

            //2. Show loading icon.
            $('.boxLoading').fadeIn(500);

            // var form = $('#signup-form');
            // var url = form.attr("action");
            // $.post(url, form.serialize(), function (data) {
            //     if (data.success) {
            //         window.location.href = '/panel/index';
            //     } else {
            //         //Signup failed
            //         toastr.warning(data.message);
            //         return false;
            //     }
            // }, "json");
        });
    };

    return {
        //main function to initiate the module
        init: function () {
            createHandler();
        }
    };
}();