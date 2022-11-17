import React from 'react';
import { Link } from 'react-router-dom-old';
import PropTypes from 'prop-types';

import styles from './ReviewSection.module.scss';

const ReviewSection = ({ fieldData, title, editLink, useH4, datatestid }) => {
  const reviewSectionInputs = (fields) => {
    return fields.map((field) => (
      <tr key={field.label}>
        <th scope="row">
          <strong>{field.label}</strong>
        </th>
        <td>{field.value}</td>
      </tr>
    ));
  };

  return (
    <div>
      {title && (
        <div>
          {!useH4 ? (
            <h3>
              {title}
              <span className={styles['edit-section-link']}>
                <Link data-testid={datatestid} to={editLink} className="usa-link">
                  Edit
                </Link>
              </span>
            </h3>
          ) : (
            <h4>
              {title}
              <span className={styles['edit-section-link']}>
                <Link data-testid={datatestid} to={editLink} className="usa-link">
                  Edit
                </Link>
              </span>
            </h4>
          )}
        </div>
      )}
      <table className={styles['review-section']}>
        <colgroup>
          <col style={{ width: '33%' }} />
          <col style={{ width: '64%' }} />
        </colgroup>
        <tbody>{reviewSectionInputs(fieldData)}</tbody>
      </table>
    </div>
  );
};

ReviewSection.defaultProps = {
  title: '',
  editLink: '',
  useH4: false,
  datatestid: '',
};

ReviewSection.propTypes = {
  fieldData: PropTypes.arrayOf(
    PropTypes.shape({
      label: PropTypes.string,
      value: PropTypes.node,
      key: PropTypes.string,
    }),
  ).isRequired,
  title: PropTypes.string,
  editLink: PropTypes.string,
  datatestid: PropTypes.string,
  useH4: PropTypes.bool,
};

export default ReviewSection;
