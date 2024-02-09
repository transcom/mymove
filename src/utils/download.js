/**
 * On REST onsuccess handler for PPM AOA packet download.
 * @param {httpResponse} response
 */
export const downloadPPMAOAPacketOnSuccessHandler = (response) => {
  // dynamically update DOM to trigger browser to display SAVE AS download file modal
  const url = window.URL.createObjectURL(new Blob([response.data]));

  const link = document.createElement('a');
  link.href = url;
  const disposition = response.headers['content-disposition'];
  const filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
  let filename = 'ppmAOAPacket.pdf';
  const matches = filenameRegex.exec(disposition);
  if (matches != null && matches[1]) {
    filename = matches[1].replace(/['"]/g, '');
  }
  link.setAttribute('download', filename);

  document.body.appendChild(link);

  // Start download
  link.click();

  // Clean up and remove the link
  link.parentNode.removeChild(link);
};

export default downloadPPMAOAPacketOnSuccessHandler;
