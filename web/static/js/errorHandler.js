export function handleAfterRequest(event) {
    const status = event.detail.xhr.status;
    if (status === 401 || status === 400) {
        showErrorModal(event.detail.xhr.response);
    }
}

export function showErrorModal(response) {
    const errorModal = document.getElementById("error-modal");
    const responseJson = JSON.parse(response);
    const errorMessage = responseJson.error;

    const errorModalTxt = document.getElementById("error-modal-message");
    errorModalTxt.textContent = errorMessage;

    errorModal.classList.remove("hidden");
}

export function closeErrorModal() {
    const errorModal = document.getElementById("error-modal");
    errorModal.classList.add("hidden");
}