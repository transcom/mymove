import { Button } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import { connect } from 'react-redux';

import styles from './AsyncPacketDownloadLink.module.scss';

import { setShowLoadingSpinner as setShowLoadingSpinnerAction } from 'store/general/actions';

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
 * @param {func} onSuccess on success response handler
 * @param {func} onFailure on failure response handler
 * @param {func} setShowLoadingSpinner used for loading spinner mask
 * @param {string} loadingMessage used for setting the loading message on spinner mask
 */
const AsyncPacketDownloadLink = ({
  id,
  label,
  asyncRetrieval,
  onSuccess,
  onFailure,
  className,
  setShowLoadingSpinner,
  loadingMessage,
}) => {
  const dataTestId = `asyncPacketDownloadLink${id}`;

  const handleClick = () => {
    setShowLoadingSpinner(true, loadingMessage);
    asyncRetrieval(id)
      .then((response) => {
        onSuccess(response);
        setShowLoadingSpinner(false, null);
      })
      .catch(() => {
        onFailure();
        setShowLoadingSpinner(false, null);
      });
  };

  return (
    <Button
      data-testid={dataTestId}
      className={className ? className : styles.downloadButtonToLink}
      onClick={handleClick}
    >
      {label}
    </Button>
  );
};

AsyncPacketDownloadLink.propTypes = {
  id: PropTypes.string.isRequired,
  label: PropTypes.string.isRequired,
  asyncRetrieval: PropTypes.func.isRequired,
  onSuccess: PropTypes.func.isRequired,
  onFailure: PropTypes.func.isRequired,
  className: PropTypes.string,
  setShowLoadingSpinner: PropTypes.func,
  loadingMessage: PropTypes.string,
};

AsyncPacketDownloadLink.defaultProps = {
  onSuccess: onPacketDownloadSuccessHandler,
  onFailure: () => {},
  setShowLoadingSpinner: () => {},
  loadingMessage: null,
};

const mapDispatchToProps = {
  setShowLoadingSpinner: setShowLoadingSpinnerAction,
};

export default connect(() => ({}), mapDispatchToProps)(AsyncPacketDownloadLink);
