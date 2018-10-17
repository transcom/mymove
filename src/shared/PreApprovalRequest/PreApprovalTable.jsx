import React from 'react';
import PropTypes from 'prop-types';
import PreApprovalRequest from 'shared/PreApprovalRequest/PreApprovalRequest.jsx';

import './PreApprovalRequest.css';

const PreApprovalTable = ({ shipment_accessorials, isActionable, onEdit, onApproval, onDelete }) => (
  <div className="accessorial-panel">
    <table cellSpacing={0}>
      <tbody>
        <tr>
          <th>Code</th>
          <th>Item</th>
          <th>Loc.</th>
          <th>Base Quantity</th>
          <th>Notes</th>
          <th>Submitted</th>
          <th>Status</th>
          <th>&nbsp;</th>
        </tr>
        {shipment_accessorials.map(row => {
          return (
            <PreApprovalRequest
              key={row.id}
              shipmentLineItem={row}
              onEdit={onEdit}
              onApproval={onApproval}
              onDelete={onDelete}
              isActionable={isActionable}
            />
          );
        })}
      </tbody>
    </table>
  </div>
);

PreApprovalTable.propTypes = {
  shipment_accessorials: PropTypes.array,
  isActionable: PropTypes.bool,
  onEdit: PropTypes.func,
  onDelete: PropTypes.func,
  onApproval: PropTypes.func,
};

export default PreApprovalTable;
