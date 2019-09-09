import React from 'react';
import { string } from 'prop-types';
import deleteButtonImg from 'shared/images/delete-doc-button.png';

const ExpenseTicketListItem = ({ id, amount, type, paymentMethod, showDelete, deleteDocumentListItem }) => {
  return (
    <div className="ticket-item">
      <div className="expense-li-item-container">
        <h4>
          {type} - ${amount}
        </h4>
        {showDelete && (
          <img alt="delete document button" onClick={() => deleteDocumentListItem(id)} src={deleteButtonImg} />
        )}
      </div>
      <div>
        {type} ({paymentMethod === 'OTHER' ? 'Not GTCC' : paymentMethod})
      </div>
    </div>
  );
};

ExpenseTicketListItem.propTypes = {
  id: string.isRequired,
  amount: string.isRequired,
  type: string.isRequired,
  paymentMethod: string.isRequired,
};

ExpenseTicketListItem.defaultProps = {
  showDelete: false,
};

export default ExpenseTicketListItem;
