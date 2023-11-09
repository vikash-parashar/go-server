// <!-- offcanvas script -->

$(document).ready(function () {
    // Initialize the offcanvas component
    var offcanvas = new bootstrap.Offcanvas(document.getElementById('offcanvasExample'));

    // Track the offcanvas state
    var offcanvasOpen = false;

    // Function to toggle the offcanvas state
    function toggleOffcanvas() {
        offcanvasOpen = !offcanvasOpen;
        if (offcanvasOpen) {
            offcanvas.show();
        } else {
            offcanvas.hide();
        }
    }

    // Event listener for mouseenter on the left side of the screen
    $(document).on('mouseenter', function (event) {
        if (event.clientX < 50) { // Adjust the value as needed
            toggleOffcanvas();
        }
    });

    // Event listener for mouseleave on the offcanvas element
    $('#offcanvasExample').on('mouseleave', function () {
        toggleOffcanvas();
    });
});




// <!-- sweet alerts -->

function showToast(icon, title) {
    const Toast = Swal.mixin({
        toast: true,
        position: 'top-end',
        showConfirmButton: false,
        timer: 3000,
        timerProgressBar: false,
        didOpen: (toast) => {
            toast.addEventListener('mouseenter', Swal.stopTimer);
            toast.addEventListener('mouseleave', Swal.resumeTimer);
        }
    });

    Toast.fire({
        icon: icon,
        title: title
    });
}
// showToast('success', 'Signed in successfully');


// <!-- handle Delete Record -->

// Add event listeners to the buttons
const deleteButtons = document.querySelectorAll(".delete-button");
deleteButtons.forEach(button => {
    button.addEventListener("click", function () {
        const id = this.getAttribute("data-id");

        Swal.fire({
            title: 'Are you sure?',
            text: "You won't be able to revert this!",
            icon: 'warning',
            showCancelButton: true,
            confirmButtonText: 'Yes, delete it!',
            cancelButtonText: 'No, cancel',
        }).then((result) => {
            if (result.isConfirmed) {
                deleteRecord(id);
            }
        });
    });
});

// Function to send a DELETE request
function deleteRecord(id) {
    fetch(`http://localhost:8080/api/v1/location-details/${id}`, {
        method: 'DELETE',
    })
        .then(response => {
            if (response.status === 200) {
                showToast('success', 'Entry deleted Successfully.');
                location.reload();
            } else {
                showToast('error', 'Try Again!');
            }
        })
        .then(data => {
            showToast("success", data.message); // Show success alert
            location.reload(); // Reload the page
        })
        .catch(error => {
            console.error("Error:", error);
            showToast("error", error.message); // Show error alert
        });
}

// Helper function to display a toast message
function showToast(type, message) {
    // Implement your toast message display logic here
    console.log(type, message);
}


// <!-- handle get data -->

const tableRows = document.querySelectorAll(".row-transition");

tableRows.forEach((row, index) => {
    row.style.opacity = "0";
    row.style.transform = "translateY(20px)";
    setTimeout(() => {
        row.style.transition = "opacity 0.2s ease, transform 0.5s ease";
        row.style.opacity = "1";
        row.style.transform = "translateY(0)";
    }, 500 * (index + 1)); // 500ms (0.5s) delay for each row
});


// <!-- handle add new data -->

document.getElementById("entryForm").addEventListener("submit", function (e) {
    e.preventDefault(); // Prevent the default form submission

    // Get form data
    const formData = new FormData(document.getElementById("entryForm"));

    // Send the form data to your backend API using fetch or another AJAX method

    var xhr = new XMLHttpRequest();
    xhr.open("POST", "http://localhost:8080/api/v1/location-details", true);
    xhr.onreadystatechange = function () {
        if (xhr.readyState === 4 && xhr.status === 200) {
            // Request completed, handle response
            var data = JSON.parse(xhr.responseText);
            console.log(xhr.responseText);

            if (data.success) {
                notify(1, data.message);
                location.reload();
            } else {
                notify(3, response.message);
            }
        }
    };
    xhr.send(formData);
    // Close the modal

    $("#createEntryModal").modal("hide");
});



// <!-- handle update existing data -->

$(document).ready(function () {
    // Handle the click event on the "Update" button for Location
    $(".update-button").on("click", function () {
        var recordId = $(this).data("id");
        $("#updateLocationId").val(recordId);
        $("#updateLocationModal").modal("show");
    });

    // Handle the form submission for Location
    $("#updateLocationSubmit").on("click", function () {
        var formData = {};
        formData.id = $("#updateLocationId").val();
        formData.serial_number = $("#updateSerialNumber").val();
        formData.device_make_model = $("#updateDeviceMakeModel").val();
        formData.model = $("#updateModel").val();
        formData.device_type = $("#updateDeviceType").val();
        formData.data_center = $("#updateDataCenter").val();
        formData.region = $("#updateRegion").val();
        formData.dc_location = $("#updateDCLocation").val();
        formData.device_location = $("#updateDeviceLocation").val();
        formData.device_row_number = $("#updateDeviceRowNumber").val();
        formData.device_rack_number = $("#updateDeviceRackNumber").val();
        formData.device_ru_number = $("#updateDeviceRUNumber").val();

        $.ajax({
            type: "PUT",
            url: "/api/v1/location-details/" + formData.id,
            contentType: "application/json", // Set the correct content type
            data: JSON.stringify(formData),
            success: function (response) {
                if (response.success) {
                    alert("Device Location Detail updated successfully");
                    $("#updateLocationModal").modal("hide");
                    location.reload(); // Reload the page after a successful update
                } else {
                    alert("Failed to update Device Location Detail");
                }
            },
            error: function (xhr, status, error) {
                // Handle errors here
                alert("An error occurred while updating: " + error);
            },
        });
    });
});




// <!-- handle filter data table -->

function filterTable() {
    var serialNumberFilter = document.getElementById("serialNumberFilter").value.toLowerCase();
    var makeFilter = document.getElementById("makeFilter").value.toLowerCase();
    var modelFilter = document.getElementById("modelFilter").value.toLowerCase();

    var rows = document.querySelectorAll("#deviceDetails tr.row-transition");

    rows.forEach(function (row) {
        var serialNumber = row.querySelector("td:nth-child(2)").textContent.toLowerCase();
        var make = row.querySelector("td:nth-child(3)").textContent.toLowerCase();
        var model = row.querySelector("td:nth-child(4)").textContent.toLowerCase();

        if (serialNumber.includes(serialNumberFilter) &&
            make.includes(makeFilter) &&
            model.includes(modelFilter)) {
            row.style.display = "";
        } else {
            row.style.display = "none";
        }
    })
}

// Add an event listener to the "Apply Filter" button
document.getElementById("applyFilterButton").addEventListener("click", function (event) {
    event.preventDefault(); // Prevent the default form submission
    filterTable(); // Apply the filter
});


// <!-- for download data -->

document.getElementById("downloadPDFButton").addEventListener("click", function () {
    // Show a dialog to let the user choose PDF or XLSX
    if (confirm("Download as PDF?")) {
        // Make a request to the PDF download API
        fetch("/api/v1/fiber-details/pdf")
            .then(response => response.blob())
            .then(data => {
                // Create a Blob containing the PDF data
                const blob = new Blob([data], { type: "application/pdf" });

                // Create a URL for the Blob
                const url = window.URL.createObjectURL(blob);

                // Create an invisible anchor element to trigger the download
                const a = document.createElement("a");
                a.style.display = "none";
                a.href = url;
                a.download = "DeviceEthernetFiberDetails.pdf";

                // Trigger the download
                document.body.appendChild(a);
                a.click();

                // Clean up
                window.URL.revokeObjectURL(url);
            })
            .catch(error => {
                console.error("Failed to download PDF:", error);
            });
    }
});

document.getElementById("downloadXLSXButton").addEventListener("click", function () {
    // Show a dialog to let the user choose PDF or XLSX
    if (confirm("Download as XLSX?")) {
        // Make a request to the XLSX download API
        fetch("/api/v1/fiber-details/excel")
            .then(response => response.blob())
            .then(data => {
                // Create a Blob containing the XLSX data
                const blob = new Blob([data], { type: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" });

                // Create a URL for the Blob
                const url = window.URL.createObjectURL(blob);

                // Create an invisible anchor element to trigger the download
                const a = document.createElement("a");
                a.style.display = "none";
                a.href = url;
                a.download = "DeviceEthernetFiberDetails.xlsx";

                // Trigger the download
                document.body.appendChild(a);
                a.click();

                // Clean up
                window.URL.revokeObjectURL(url);
            })
            .catch(error => {
                console.error("Failed to download XLSX:", error);
            });
    }
});



// <!-- get current user details -->

function notify(type, msg) {
    notie.alert({
        type: type,
        text: msg,
        position: 'top'
    })
}
function capitalize(str) {
    return str.slice(0).toUpperCase();
}


// Function to fetch and display user details when the page loads
function fetchCurrentUserDetails() {
    fetch('/api/v1/get-current-user', {
        method: 'GET',
        headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
        }
    })
        .then(response => {
            console.log(response)
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(data => {
            document.getElementById('firstNamePlaceholder').innerHTML = capitalize(data.user.first_name);
            // document.getElementById('lastNamePlaceholder').innerHTML = capitalize(data.user.last_name);
            document.getElementById('phonePlaceholder').innerHTML = capitalize(data.user.phone);
            // document.getElementById('emailPlaceholder').innerHTML = capitalize(data.user.email);
            document.getElementById('rolePlaceholder').innerHTML = capitalize(data.user.role);
        })

        .catch(error => {
            console.error('Error fetching user details:', error);
        });
}

// Add an event listener to fetch user details when the page loads
document.addEventListener('DOMContentLoaded', fetchCurrentUserDetails);


// <!-- handle logout -->

// Function to handle logout
function handleLogout() {
    fetch("/logout", {
        method: "POST",
    })
        .then(response => {
            if (!response.ok) {
                showToast("warning", "Try Again !");
                console.error("Network response was not ok");
            }
            return response.json();
        })
        .then(data => {
            if (data.success) {
                showToast("success", "Logout Success.");
                window.location.href = "/"; // Redirect to the root path
            } else {
                showToast("warning", "Try Again !");
            }
        })
        .catch(error => {
            console.error("An error occurred while logging out:", error);
            showToast("error", "Try Again !");
        });
}

// Attach the handleLogout function to the logout button's click event
const logoutButton = document.getElementById("logout-button");
logoutButton.addEventListener("click", handleLogout);
