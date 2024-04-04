import React from 'react';
import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';

import styles from './AsyncPacketDownloadLink.module.scss';

export const onPacketDownloadSuccessHandler = (response) => {
  // dynamically update DOM to trigger browser to display SAVE AS download file modal
  const contentType = response.headers['content-type'];
  const url = window.URL.createObjectURL(
    new Blob([response.data], {
      type: contentType,
    }),
  );

  const link = document.createElement('a');
  link.href = url;
  const disposition = response.headers['content-disposition'];
  const filenameRegex = /filename[^;=\n]*=((['"]).*?\2|[^;\n]*)/;
  let strtime = new Date().toISOString();
  // we are expecting PDF
  let filename = `packet-${strtime}.pdf`;
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

/**
 * Shared component to render download links for AOA/Payment Packet packets.
 * @param {string} id uuid to download
 * @param {string} label link text
 * @param {Promise} asyncRetrieval asynch document retrieval
 * @param {func} onSucccess on success response handler
 * @param {func} onFailure on failure response handler
 */
const AsyncPacketDownloadLink = ({ id, label, asyncRetrieval, onSucccess, onFailure, className }) => {
  const dataTestId = `asyncPacketDownloadLink${id}`;

  return (
    <Button
      data-testid={dataTestId}
      className={className ? className : styles.downloadButtonToLink}
      onClick={() =>
        asyncRetrieval(id)
          .then((response) => {
            onSucccess(response);
          })
          .catch(() => {
            onFailure();
          })
      }
    >
      {label}
    </Button>
  );
};

AsyncPacketDownloadLink.propTypes = {
  id: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  asyncRetrieval: PropTypes.func.isRequired,
  onSucccess: PropTypes.func.isRequired,
  onFailure: PropTypes.func.isRequired,
  className: PropTypes.string,
};

AsyncPacketDownloadLink.defaultProps = {
  onSucccess: onPacketDownloadSuccessHandler,
  onFailure: () => {},
};

export default AsyncPacketDownloadLink;
