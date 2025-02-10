$(document).ready(function () {
    $('[data-toggle="tooltip"]').tooltip();
    $("#apiResetForm").submit(function (e) {
        api.reset()
            .success(function (response) {
                user.api_key = response.data
                    
                $("#api_key").val(user.api_key)
            })
            .error(function (data) {
                errorFlash(data.message)
            })
        return false
    })
    $("#settingsForm").submit(function (e) {
        $.post("/settings", $(this).serialize())
            .done(function (data) {
                successFlash(data.message)
            })
            .fail(function (data) {
                errorFlash(data.responseJSON.message)
            })
        return false
    })
    //$("#imapForm").submit(function (e) {
    $("#savesettings").click(function() {
        var imapSettings = {}
        imapSettings.host = $("#imaphost").val()
        imapSettings.port = $("#imapport").val()
        imapSettings.username = $("#imapusername").val()
        imapSettings.password = $("#imappassword").val()
        imapSettings.enabled = $('#use_imap').prop('checked')
        imapSettings.tls = $('#use_tls').prop('checked')

        //Advanced settings
        imapSettings.folder = $("#folder").val()
        imapSettings.imap_freq = $("#imapfreq").val()
        imapSettings.restrict_domain = $("#restrictdomain").val()
        imapSettings.ignore_cert_errors = $('#ignorecerterrors').prop('checked')
        imapSettings.delete_reported_campaign_email = $('#deletecampaign').prop('checked')
        
        //To avoid unmarshalling error in controllers/api/imap.go. It would fail gracefully, but with a generic error.
        if (imapSettings.host == ""){
            errorFlash("No IMAP Host specified")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (imapSettings.port == ""){
            errorFlash("No IMAP Port specified")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (isNaN(imapSettings.port) || imapSettings.port <1 || imapSettings.port > 65535  ){ 
            errorFlash("Invalid IMAP Port")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (imapSettings.imap_freq == ""){
            imapSettings.imap_freq = "60"
        }

        api.IMAP.post(imapSettings).done(function (data) {
                if (data.success == true) {
                    successFlashFade("Successfully updated IMAP settings.", 2)
                } else {
                    errorFlash("Unable to update IMAP settings.")
                }
            })
            .success(function (data){
                loadIMAPSettings()
            })
            .fail(function (data) {
                errorFlash(data.responseJSON.message)
            })
            .always(function (data){
                document.body.scrollTop = 0;
                document.documentElement.scrollTop = 0;
            })
        
        return false
    })

    $("#validateimap").click(function() {

        // Query validate imap server endpoint
        var server = {}
        server.host = $("#imaphost").val()
        server.port = $("#imapport").val()
        server.username = $("#imapusername").val()
        server.password = $("#imappassword").val()
        server.tls = $('#use_tls').prop('checked')
        server.ignore_cert_errors = $('#ignorecerterrors').prop('checked')

        //To avoid unmarshalling error in controllers/api/imap.go. It would fail gracefully, but with a generic error. 
        if (server.host == ""){
            errorFlash("No IMAP Host specified")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (server.port == ""){
            errorFlash("No IMAP Port specified")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }
        if (isNaN(server.port) || server.port <1 || server.port > 65535  ){
            errorFlash("Invalid IMAP Port")
            document.body.scrollTop = 0;
            document.documentElement.scrollTop = 0;
            return false
        }

        var oldHTML = $("#validateimap").html();
        // Disable inputs and change button text
        $("#imaphost").attr("disabled", true);
        $("#imapport").attr("disabled", true);
        $("#imapusername").attr("disabled", true);
        $("#imappassword").attr("disabled", true);
        $("#use_imap").attr("disabled", true);
        $("#use_tls").attr("disabled", true);
        $('#ignorecerterrors').attr("disabled", true);
        $("#folder").attr("disabled", true);
        $("#restrictdomain").attr("disabled", true);
        $('#deletecampaign').attr("disabled", true);
        $('#lastlogin').attr("disabled", true);
        $('#imapfreq').attr("disabled", true);
        $("#validateimap").attr("disabled", true);  
        $("#validateimap").html("<i class='fa fa-circle-o-notch fa-spin'></i> Testing...");
        
        api.IMAP.validate(server).done(function(data) {
            if (data.success == true) {
                Swal.fire({
                    title: "Success",
                    html: "Logged into <b>" + escapeHtml($("#imaphost").val()) + "</b>",
                    type: "success",
                })
            } else {
                Swal.fire({
                    title: "Failed!",
                    html: "Unable to login to <b>" + escapeHtml($("#imaphost").val()) + "</b>.",
                    type: "error",
                    showCancelButton: true,
                    cancelButtonText: "Close",
                    confirmButtonText: "More Info",
                    confirmButtonColor: "#428bca",
                    allowOutsideClick: false,
                }).then(function(result) {
                    if (result.value) {
                        Swal.fire({
                            title: "Error:",
                            text: data.message,
                        })
                    }
                  })
            }
            
          })
          .fail(function() {
            Swal.fire({
                title: "Failed!",
                text: "An unecpected error occured.",
                type: "error",
            })
          })
          .always(function() {
            //Re-enable inputs and change button text
            $("#imaphost").attr("disabled", false);
            $("#imapport").attr("disabled", false);
            $("#imapusername").attr("disabled", false);
            $("#imappassword").attr("disabled", false);
            $("#use_imap").attr("disabled", false);
            $("#use_tls").attr("disabled", false);
            $('#ignorecerterrors').attr("disabled", false);
            $("#folder").attr("disabled", false);
            $("#restrictdomain").attr("disabled", false);
            $('#deletecampaign').attr("disabled", false);
            $('#lastlogin').attr("disabled", false);
            $('#imapfreq').attr("disabled", false);
            $("#validateimap").attr("disabled", false);
            $("#validateimap").html(oldHTML);

          });

      }); //end testclick

    $("#reporttab").click(function() {
        loadIMAPSettings()
    })

    $("#advanced").click(function() {
        $("#advancedarea").toggle();
    })

    $("#saveclientdata").click(function() {
        let client = {
            name: $("#client_name").val(),
            email: $("#client_email").val(),
            monitor_url: $("#client_monitor_url").val(),
            monitor_password: $("#client_monitor_password").val(),
            apolo_api_key: $("#client_api_key").val()
        };
    
        api.client.post(client).done(function(data) {
            if (data.success) {
                successFlash(data.message);
                loadClientData();
                loadClientHistory();
            } else {
                errorFlash(data.message);
            }
        }).fail(function() {
            errorFlash("Error saving client data.");
        });
    
        return false;
    });
    
    
    $("#cancelclient").click(function() {
	loadClientData()
    })

    $("#saveqrsettings").click(function() {
        qr = {}
        qr.qr_size = parseInt($("#qr_size").val())
        qr.qr_pixels = $("#qr_pixels").val()
        qr.qr_background = $("#qr_background").val()

        api.QR.post(qr).done(function(data) {
            if (data.success == true) {
                successFlash(data.message)
            } else {
                errorFlash(data.message)
            }
        })
        return false
    })

    $("#cancelqr").click(function() {
        loadQRConfigs()
    })

    function loadIMAPSettings(){
        api.IMAP.get()
        .success(function (imap) {
            if (imap.length == 0){
                $('#lastlogindiv').hide()
            } else {
                imap = imap[0]
                if (imap.enabled == false){
                    $('#lastlogindiv').hide()
                } else {
                    $('#lastlogindiv').show()
                }
                $("#imapusername").val(imap.username)
                $("#imaphost").val(imap.host)
                $("#imapport").val(imap.port)
                $("#imappassword").val(imap.password)
                $('#use_tls').prop('checked', imap.tls)
                $('#ignorecerterrors').prop('checked', imap.ignore_cert_errors)
                $('#use_imap').prop('checked', imap.enabled)
                $("#folder").val(imap.folder)
                $("#restrictdomain").val(imap.restrict_domain)
                $('#deletecampaign').prop('checked', imap.delete_reported_campaign_email)
                $('#lastloginraw').val(imap.last_login)
                $('#lastlogin').val(moment.utc(imap.last_login).fromNow())
                $('#imapfreq').val(imap.imap_freq)
            }  

        })
        .error(function () {
            errorFlash("Error fetching IMAP settings")
        })
    }

    function loadClientData() {
        api.client.get()
            .success(function(client) {
                if (!client || Object.keys(client).length === 0) {
                    errorFlash("No clients registered.");
                    return;
                }
    
                $("#client_name").val(client.name);
                $("#client_name_modal").val(client.name);
                $("#client_email").val(client.email);
                $("#client_email_modal").val(client.email);
                $("#client_monitor_url").val(client.monitor_url);
                $("#client_monitor_url_modal").val(client.monitor_url);
                $("#client_monitor_password").val(client.monitor_password);
                $("#client_monitor_password_modal").val(client.monitor_password);
                $("#client_api_key").val(client.apolo_api_key);
                $("#client_api_key_modal").val(client.apolo_api_key);
    
                let isDisabled = !client.email || !client.apolo_api_key;
                $("#send_monitor_modal").prop("disabled", isDisabled);
    
                $("#send_status").find("span").remove();
    
                let statusHtml = '';
    
                if (!client.send_date || client.send_date === "Not sent") {
                    statusHtml += `<span class="badge" style="background-color: rgb(224, 51, 51); color:white;">
                                    <i class="bi bi-check-circle"></i> Not sent
                                   </span> `;
                }
                if (!client.email) {
                    statusHtml += `<span class="badge" style="background-color: rgb(199, 201, 79); color:white;">
                                    <i class="bi bi-check-circle"></i> Not registered client email
                                   </span> `;
                }
                if (!client.monitor_password) {
                    statusHtml += `<span class="badge" style="background-color: rgb(255, 165, 0); color:white;">
                                    <i class="bi bi-exclamation-triangle"></i> Not registered monitor password
                                   </span>`;
                }
                if (!client.apolo_api_key) {
                    statusHtml += `<span class="badge" style="background-color: rgb(199, 201, 79); color:white;">
                                    <i class="bi bi-check-circle"></i> Not registered api key
                                   </span>`;
                }
    
                $("#send_status").append(statusHtml);
            })
            .error(function() {
                errorFlash("Error fetching client data.");
            });
    }
    
    
    function loadMonitorUrl() {
        const domainName = window.location.hostname.split('.').slice(-2).join('.');
        /* Receive and load unexposed URL*/ 
        /* Load from server*/ 
        const monitorUrl = `https:/${domainName}/monitor`;
        
        const inputURL = document.getElementById('client_monitor_url');
        if (inputURL) {
            inputURL.value = monitorUrl;
        }
        const inputURLsending = document.getElementById('sending_phishing_monitor_url');
        if (inputURLsending) {
            inputURLsending.value = monitorUrl;
        }
        
        const linkElement = document.getElementById('go_to_panel_monitor');
        if (linkElement) {
            linkElement.href = monitorUrl;
        }
    }

    $(document).on("click", ".lock-field", function() {
        var input = $($(this).data("target"));
        var lockButton = $(this);

        if (input.prop("readonly")) {
            input.prop("readonly", false);
            lockButton.html('<i class="fa fa-lock"></i>');
        } else {
            input.prop("readonly", true);
            lockButton.html('<i class="fa fa-unlock"></i>');
        }
    });

    $('.copy-button').on('click', function() {
        var targetId = $(this).data('target');
        var $inputElement = $(targetId);
    
        if ($inputElement.length) {
            navigator.clipboard.writeText($inputElement.val())
                .catch(err => {
                    console.error('Failed to copy text: ', err);
                });
        }
    });

    $(document).on("click", ".reveal-password", function() {
        var input = $($(this).data("target"));
        var RevealButton = $(this);
    
        if (input.attr("type") === "password") {
            input.attr("type", "text");
            RevealButton.html('<i class="fa fa-eye-slash"></i>');
        } else {
            input.attr("type", "password");
            RevealButton.html('<i class="fa fa-eye"></i>');
        }
    });

    function loadClientHistory() {
        api.client_history.get()
            .done(function(history) {
                if (!Array.isArray(history)) {
                    console.error("Expected an array but got:", history);
                    return;
                }
    
                $("#clientTable tbody").empty();
    
                history.forEach(function (record) {
                    let email = record.email ? record.email : 'Not provided';
                    let createdAt = record.created_at ? new Date(record.created_at).toLocaleString() : 'Unknown';
                    let sendDate = record.send_date ? new Date(record.send_date).toLocaleString() : 'Not sent';
                    let changeDate = record.change_date ? new Date(record.change_date).toLocaleString() : 'Unknown';
    
                    let row = `<tr>
                        <td>${record.name}</td>
                        <td>${record.monitor_url}</td>
                        <td><span class="hidden-password" data-password="${record.monitor_password}">**********</span></td>
                        <td><span class="hidden-password" data-password="${record.apolo_api_key}">**********</span></td>
                        <td>${createdAt}</td>
                        <td>${email}</td>
                        <td>${sendDate}</td>
                        <td>${changeDate}</td>
                    </tr>`;
    
                    $("#clientTable tbody").append(row);
                });
            })
            .fail(function() {
                errorFlash("Error fetching client history data");
            });
    }
    
    
    $(document).on("click", "#reveal_password_from_client_table", function() {
        let button = $(this);
        let isHidden = button.find("i").hasClass("fa-eye");
        
        $(".hidden-password").each(function() {
            let span = $(this);
            if (isHidden) {
                span.text(span.data("password"));
            } else {
                span.text("**********");
            }
        });
    
        if (isHidden) {
            button.html('<i class="fa fa-eye-slash"></i> Hide passwords');
        } else {
            button.html('<i class="fa fa-eye"></i> Reveal passwords');
        }
    });
    
    function loadQRConfigs() {
        api.QR.get()
        .success(function (qr) {
            $("#qr_size").val(qr.qr_size)
            $("#qr_pixels").val(qr.qr_pixels)
            $("#qr_background").val(qr.qr_background)
        })
    }

    var use_map = localStorage.getItem('gophish.use_map')
    $("#use_map").prop('checked', JSON.parse(use_map))
    $("#use_map").on('change', function () {
        localStorage.setItem('gophish.use_map', JSON.stringify(this.checked))
    })

    loadIMAPSettings()
    loadQRConfigs()
    loadClientData()
    loadMonitorUrl()
    loadClientHistory()
})
