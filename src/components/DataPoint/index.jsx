import React from 'react';
import propTypes from 'prop-types';

const DataPoint = ({ header, body, custClass }) => (
  <table className={`table--data-point ${custClass}`}>
    <thead className="table--small">
      <tr>
        <th>{header}</th>
      </tr>
    </thead>
    <tbody>
      <tr>
        <td>{body}</td>
      </tr>
    </tbody>
  </table>
);

DataPoint.propTypes = {
  header: propTypes.string.isRequired,
  body: propTypes.element.isRequired,
  custClass: propTypes.string,
};

DataPoint.defaultProps = {
  custClass: '',
};

export default DataPoint;
