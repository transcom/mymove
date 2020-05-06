import React from 'react';
import propTypes from 'prop-types';

const DataPoint = ({ header, body }) => (
  <table className="table--data-point">
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
  header: propTypes.string,
  body: propTypes.element,
};

export default DataPoint;
