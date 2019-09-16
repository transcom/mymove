import React, { Component } from 'react';
import { string } from 'prop-types';
import deleteButtonImg from 'shared/images/delete-doc-button.png';
import AlertWithDeleteConfirmation from 'shared/AlertWithDeleteConfirmation';

class ExpenseTicketListItem extends Component {
  state = {
    showDeleteConfirmation: false,
  };

  toggleShowConfirmation = () => {
    const { showDeleteConfirmation } = this.state;
    this.setState({ showDeleteConfirmation: !showDeleteConfirmation });
  };

  render() {
    const { id, amount, type, paymentMethod, showDelete, deleteDocumentListItem } = this.props;
    const { showDeleteConfirmation } = this.state;
    return (
      <div className="ticket-item" style={{ display: 'flex' }}>
        <div style={{ flex: 1 }}>
          <div className="expense-li-item-container">
            <h4>
              {type} - ${amount}
            </h4>
            {showDelete && (
              <img
                alt="delete document button"
                data-cy="delete-ticket"
                onClick={this.toggleShowConfirmation}
                src={deleteButtonImg}
              />
            )}
          </div>
          <div>
            {type} ({paymentMethod === 'OTHER' ? 'Not GTCC' : paymentMethod})
          </div>
          {showDeleteConfirmation && (
            <AlertWithDeleteConfirmation
              heading="Delete this document?"
              message="This action cannot be undone."
              deleteActionHandler={() => deleteDocumentListItem(id)}
              cancelActionHandler={this.toggleShowConfirmation}
              type="expense-ticket-list-alert"
            ></AlertWithDeleteConfirmation>
          )}
        </div>
      </div>
    );
  }
}

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
