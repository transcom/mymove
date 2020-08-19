import React from 'react';
import { Link } from 'react-router-dom';
import PropTypes from 'prop-types';

import 'scenes/Review/Review.css';

const ReviewSection = ({ fieldData, title, editLink }) => {
  const reviewSectionInputs = (fields) => {
    return fields.map((field) => (
      <tr key={field.label}>
        <th scope="row">{field.label}</th>
        <td>{field.value}</td>
      </tr>
    ));
  };

  return (
    <div>
      {title && (
        <div>
          <h2>
            {title}
            <span className="edit-section-link">
              <Link to={editLink} className="usa-link">
                Edit
              </Link>
            </span>
          </h2>
        </div>
      )}
      <table className="review-section">
        <colgroup>
          <col style={{ width: '25%' }} />
          <col style={{ width: '75%' }} />
        </colgroup>
        <tbody>{reviewSectionInputs(fieldData)}</tbody>
      </table>
    </div>
  );
};

ReviewSection.defaultProps = {
  title: '',
  editLink: '',
};

ReviewSection.propTypes = {
  fieldData: PropTypes.arrayOf(
    PropTypes.shape({
      label: PropTypes.string,
      value: PropTypes.string,
      key: PropTypes.string,
    }),
  ).isRequired,
  title: PropTypes.string,
  editLink: PropTypes.string,
};

export default ReviewSection;
