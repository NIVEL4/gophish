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

    // Avoid sending the same client.
    let ClientData = {};

    function updateClientData() {
        $("#client_name, #client_email, #client_monitor_url, #client_monitor_password, #client_api_key").each(function () {
            ClientData[this.id] = this.value;
        });
    }

    function checkChanges() {
        if (Object.keys(ClientData).length === 0) {
            $("#saveclientdata").prop("disabled", false);
            return;
        }
    
        let hasChanges = Object.keys(ClientData).some(id => $("#" + id).val() !== ClientData[id]);
        $("#saveclientdata").prop("disabled", !hasChanges);
    }
    

    $("#client_name, #client_email, #client_monitor_url, #client_monitor_password, #client_api_key")
        .on("input change paste", checkChanges);

    $("#saveclientdata").prop("disabled", true);

    function saveClientData(specialistName = null, sendMethod = null, sentDate = null) {
        let client = {
            name: $("#client_name").val(),
            email: $("#client_email").val(),
            monitor_url: $("#client_monitor_url").val(),
            monitor_password: $("#client_monitor_password").val(),
            apolo_api_key: $("#client_api_key").val(),
        };
    
        if (specialistName !== null) {
            client.sent_by = specialistName;
        }
        if (sendMethod !== null) {
            client.send_method = sendMethod;
        }
        if (sentDate !== null) {
            client.sent_date = sentDate;
        }
    
        console.log("client save:", client);
    
        api.client.post(client).done(function(data) {
            if (data.success) {
                successFlash(data.message);
                loadClientData();
                loadClientHistory();
                loadClientSentMails();
                updateClientData();
                $("#saveclientdata").prop("disabled", true);
    
                $("#client_name, #client_email, #client_monitor_url, #client_monitor_password, #client_api_key").prop("readonly", true);
    
                $(".lock-field").each(function() {
                    $(this).html('<i class="fa fa-lock"></i>');
                });
            } else {
                errorFlash(data.message);
            }
        }).fail(function() {
            errorFlash("Error saving client data.");
        });
    
        return false;
    }
      
    function loadClientData() {
        api.client.get()
            .success(function(client) {
                if (!client || Object.keys(client).length === 0 || client.name === undefined) {
                    errorFlash("No client registered.");
                    $("#saveclientdata").prop("disabled", false);
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
    
                updateClientData();
    
                $("#saveclientdata").prop("disabled", true);
    
                let isDisabled = !client.email || !client.apolo_api_key;
                $("#send_monitor_modal").prop("disabled", isDisabled);
    
                $("#send_status").find("span").remove();
    
                let statusHtml = '';
    
                if (client.sent_by === 'Not sent' && client.send_method === 'None') {
                    statusHtml += `<span class="badge" style="background-color: rgb(224, 51, 51); color:white; font-size: 1.5rem; padding: 0.5rem 0.75rem;">
                                    <i class="bi bi-check-circle"></i> Not sent
                                   </span> `;
                    $("#sent_mails_to_client thead").css("background-color", "#b40909");
                } else {
                    statusHtml += `<span class="badge" style="background-color: rgb(105, 207, 84); color:white; font-size: 1.5rem; padding: 0.5rem 0.75rem;">
                                    <i class="bi bi-check-circle"></i> Sent
                                   </span>`;
                    $("#sent_mails_to_client thead").css("background-color", "#06a736");
                }
                if (!client.email) {
                    statusHtml += `<span class="badge" style="background-color: rgb(214, 211, 28); color:white; font-size: 1.5rem; padding: 0.5rem 0.75rem;">
                                    <i class="bi bi-check-circle"></i> Not registered client email
                                   </span> `;
                }
                if (!client.monitor_password) {
                    statusHtml += `<span class="badge" style="background-color: rgb(214, 211, 28); color:white; font-size: 1.5rem; padding: 0.5rem 0.75rem;">
                                    <i class="bi bi-exclamation-triangle"></i> Not registered monitor password
                                   </span>`;
                }
                if (!client.apolo_api_key) {
                    statusHtml += `<span class="badge" style="background-color: rgb(214, 211, 28); color:white; font-size: 1.5rem; padding: 0.5rem 0.75rem;">
                                    <i class="bi bi-check-circle"></i> Not registered api key
                                   </span>`;
                }
    
                const inputURL = document.getElementById('client_monitor_url');
                if (inputURL) {
                    inputURL.value = client.monitor_url;
                }
                const inputURLsending = document.getElementById('sending_phishing_monitor_url');
                if (inputURLsending) {
                    inputURLsending.value = client.monitor_url;
                }
    
                const linkElement = document.getElementById('go_to_panel_monitor');
                if (linkElement) {
                    linkElement.href = client.monitor_url;
                }
    
                $("#send_status").append(statusHtml);
            })
            .error(function() {
                errorFlash("Error fetching client data.");
                $("#saveclientdata").prop("disabled", false);
            });
    }    


    $("#send_monitor_via_gophish").click(function () {
        let specialistName = $("#specialist_name_modal").val();
        let sendMethod = "Gophish";
        let sentDate = new Date().toISOString();
        let emailTemplate = document.getElementById("email-template").value;
    
        let emailData = {
            client_name: $("#client_name").val(),
            client_email: $("#client_email").val(),
            client_monitor_url: $("#client_monitor_url").val(),
            client_monitor_password: $("#client_monitor_password").val(),
            specialist_name: specialistName,
            smtp_profile: parseInt($("#profile").val()),
            sent_date: sentDate,
            send_method: sendMethod,
            email_template: emailTemplate
        };
    
        api.client.send_mail_gophish(emailData)
            .done(function(response) {
                console.log("Client Gophish Mail Success: " + response);
                saveClientData(specialistName, sendMethod, sentDate);
            })
            .fail(function() {
                console.log("Client Gophish Mail Fail" + response);
            });
    });
    
    $("#send_monitor_via_apolo").click(function () {
        let specialistName = $("#specialist_name_modal").val();
        let sentDate = new Date().toISOString();
        let sendMethod = "Apolo";

        let emailData = {
            client_name: $("#client_name").val(),
            client_email: $("#client_email").val(),
            client_monitor_url: $("#client_monitor_url").val(),
            client_monitor_password: $("#client_monitor_password").val(),
            client_api_key: $("#client_api_key").val(),
            specialist_name: specialistName,
            sent_date: sentDate,
            send_method: sendMethod
        };
    
        api.client.send_mail_apolo(emailData)
            .done(function(response) {
                console.log("Client Apolo Mail Success: ", response);
                saveClientData(specialistName, sendMethod, sentDate);
            })
            .fail(function(xhr) {
                console.log("Client Apolo Mail Fail: ", xhr);
            });
    });

    $("#saveclientdata").click(function() {
        saveClientData()
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
    
    $(document).on("click", ".lock-field", function() {
        var input = $($(this).data("target"));
        var lockButton = $(this);

        if (input.prop("readonly")) {
            input.prop("readonly", false);
            lockButton.html('<i class="fa fa-unlock"></i>');
        } else {
            input.prop("readonly", true);
            lockButton.html('<i class="fa fa-lock"></i>');
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

    api.SMTP.get()
        .success(function (profiles) {
            if (profiles.length == 0) {
                modalError("No profiles found!")
                return false
            } else {
                var profile_s2 = $.map(profiles, function (obj) {
                    obj.text = obj.name
                    return obj
                });
                var profile_select = $("#profile.form-control")
                profile_select.select2({
                    placeholder: "Select a Sending Profile",
                    data: profile_s2,
                }).select2("val", profile_s2[0]);
                if (profiles.length === 1) {
                    profile_select.val(profile_s2[0].id)
                    profile_select.trigger('change.select2')
                }
            }
        });
    
    function loadClientHistory() {
        api.client.history()
            .done(function(history) {
                if (!Array.isArray(history)) {
                    console.error("Expected an array but got:", history);
                    return;
                }
    
                $("#clientTable tbody").empty();
    
                history.forEach(function (record) {
                    let email = record.email || 'Not provided';
                    let createdAt = record.created_at ? new Date(record.created_at).toLocaleString() : 'Unknown';
                    let sentDate = (record.sent_by === 'Not sent' && record.send_method === 'None') 
                        ? 'Not sent' 
                        : (record.sent_date ? new Date(record.sent_date).toLocaleString() : 'Not sent');
                    let sentBy = record.sent_by;
                    let sendMethod = record.send_method;
                    
                    let row = `<tr>
                        <td>${record.name}</td>
                        <td>${record.monitor_url}</td>
                        <td><span class="hidden-password" data-password="${record.monitor_password}">**********</span></td>
                        <td><span class="hidden-password" data-password="${record.apolo_api_key}">**********</span></td>
                        <td>${createdAt}</td>
                        <td>${email}</td>
                        <td>${sentDate}</td>
                        <td>${sentBy}</td>
                        <td>${sendMethod}</td>
                    </tr>`;
    
                    $("#clientTable tbody").append(row);
                });
            })
            .fail(function() {
                errorFlash("Error fetching client history data");
            });
    }
    
    function loadClientSentMails() {
        api.client.history()
            .done(function(history) {
                if (!Array.isArray(history)) {
                    console.error("Expected an array but got:", history);
                    return;
                }
    
                $("#sent_mails_to_client tbody").empty();
    
                history.forEach(function(record) {
                    if (record.sent_date && record.sent_by !== 'Not sent' && record.send_method !== 'None') {
                        let email = record.email;
                        let sentDate = new Date(record.sent_date).toLocaleString();
                        let sentBy = record.sent_by;
                        let sendMethod = record.send_method;

                        let row = `<tr>
                            <td>${record.name}</td>
                            <td>${email}</td>
                            <td>${sentDate}</td>
                            <td>${sentBy}</td>
                            <td>${sendMethod}</td>
                        </tr>`;
    
                        $("#sent_mails_to_client tbody").append(row);
                    }
                });
            })
            .fail(function() {
                errorFlash("Error al obtener los datos de los correos enviados");
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
    loadClientSentMails()
    loadClientHistory()
})
