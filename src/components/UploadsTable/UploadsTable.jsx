import React from 'react';
import PropTypes from 'prop-types';
import { Link } from 'react-router-dom';
import bytes from 'bytes';
import dayjs from 'dayjs';
import localizedFormat from 'dayjs/plugin/localizedFormat';
import { Button } from '@trussworks/react-uswds';

import { UPLOAD_SCAN_STATUS } from 'shared/constants';

dayjs.extend(localizedFormat);

const UploadsTable = ({ uploads, onDelete }) => {
  const getUploadUrl = (upload) => {
    switch (upload.status) {
      case UPLOAD_SCAN_STATUS.INFECTED:
        return (
          <>
            <Link to="/infected-upload" className="usa-link">
              {upload.filename}
            </Link>
          </>
        );
      case UPLOAD_SCAN_STATUS.PROCESSING:
        return (
          <>
            <Link to="/processing-upload" className="usa-link">
              {upload.filename}
            </Link>
          </>
        );
      default:
        return (
          <>
            <a href={upload.url} target="_blank" rel="noopener noreferrer" className="usa-link">
              {upload.filename}
            </a>
          </>
        );
    }
  };

  return (
    <table>
      <thead>
        <tr>
          <th>Name</th>
          <th>Uploaded</th>
          <th>Size</th>
          <th>Delete</th>
        </tr>
      </thead>
      <tbody>
        {uploads.map((upload) => (
          <tr key={upload.id} className="vertical-align text-top">
            <td className="maxw-card" style={{ overflowWrap: 'break-word', wordWrap: 'break-word' }}>
              {getUploadUrl(upload)}
            </td>
            <td>{dayjs(upload.created_at).format('LLL')}</td>
            <td>{bytes(upload.bytes)}</td>
            <td>
              <Button type="button" unstyled onClick={() => onDelete(upload.id)}>
                Delete
              </Button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
};

UploadsTable.propTypes = {
  uploads: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      created_at: PropTypes.string.isRequired,
      bytes: PropTypes.number.isRequired,
      url: PropTypes.string.isRequired,
      filename: PropTypes.string.isRequired,
    }),
  ).isRequired,
  onDelete: PropTypes.func.isRequired,
};

export default UploadsTable;
