import React from 'react';

const DataPair = () => (
  <div className="table--data-pair">
    <table className="table--data-point">
      <thead className="table--small">
        <tr>
          <th>Customer requested pick up date</th>
          <th>Scheduled pick up date</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>Thursday, 26 Mar 2020</td>
          <td>Friday, 27 Mar 2020</td>
        </tr>
      </tbody>
    </table>
  </div>
);

export default DataPair;
