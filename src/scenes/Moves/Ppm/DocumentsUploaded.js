import React, { Component } from 'react';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { connect } from 'react-redux';
import WeightTicketListItem from './PaymentReview/WeightTicketListItem';
import ExpenseTicketListItem from './PaymentReview/ExpenseTicketListItem';
import { formatExpenseDocs } from './utility';
import './PaymentReview/PaymentReview.css';

export class DocumentsUploaded extends Component {
  state = {
    showDocs: true,
  };

  defaultProps = {
    showLinks: false,
  };

  componentDidMount() {
    const { moveId } = this.props;
    this.props.getMoveDocumentsForMove(moveId);
  }

  toggleShowDocs = () => {
    this.setState({ showDocs: !this.state.showDocs });
  };

  createHeaderMessage = documentLength => {
    return (
      <div style={{ marginBottom: 5, paddingTop: 5 }}>
        {documentLength} document{documentLength > 1 ? 's' : ''} added{' '}
        <a style={{ paddingLeft: '1em' }} onClick={this.toggleShowDocs}>
          Show
        </a>
      </div>
    );
  };

  render() {
    const { showDocs } = this.state;
    const { expenseDocs, weightTicketDocs, moveId, showLinks } = this.props;
    const totalDocs = expenseDocs.length + weightTicketDocs.length;
    if (totalDocs === 0) {
      return null;
    }
    return (
      <>
        <div data-cy="documents-uploaded">
          <div className="doc-summary-container">
            <h3>Document summary - {weightTicketDocs.length + expenseDocs.length} total</h3>
            {showDocs && (
              <>
                <h4>{weightTicketDocs.length} sets of weight tickets</h4>
                <div className="tickets">
                  {weightTicketDocs.map((ticket, index) => (
                    <WeightTicketListItem key={ticket.id} num={index} {...ticket} />
                  ))}
                </div>
                {showLinks && (
                  <Link data-cy="weight-ticket-link" to={`/moves/${moveId}/ppm-weight-ticket`}>
                    <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add weight ticket
                  </Link>
                )}
                <hr id="doc-summary-separator" />
                <h4>
                  {expenseDocs.length} expense{expenseDocs.length > 1 ? 's' : ''}
                </h4>
                <div className="tickets">
                  {formatExpenseDocs(expenseDocs).map(expense => (
                    <ExpenseTicketListItem key={expense.id} {...expense} />
                  ))}
                </div>
                {showLinks && (
                  <div className="add-expense-link">
                    <Link data-cy="expense-link" to={`/moves/${moveId}/ppm-expenses`}>
                      <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} /> Add expense
                    </Link>
                  </div>
                )}
              </>
            )}
          </div>
        </div>
      </>
    );
  }
}

function mapStateToProps(state, { moveId }) {
  return {
    moveId,
    expenseDocs: selectPPMCloseoutDocumentsForMove(state, moveId, ['EXPENSE']),
    weightTicketDocs: selectPPMCloseoutDocumentsForMove(state, moveId, ['WEIGHT_TICKET_SET']),
  };
}

const mapDispatchToProps = {
  selectPPMCloseoutDocumentsForMove,
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(DocumentsUploaded);
