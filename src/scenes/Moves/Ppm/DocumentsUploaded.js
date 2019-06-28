import React, { Component } from 'react';
import { bool } from 'prop-types';
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
    showDocs: false,
  };

  defaultProps = {
    showLinks: false,
    showToggleDocs: false,
    inReviewPage: false,
  };

  componentDidMount() {
    const { moveId } = this.props;
    this.props.getMoveDocumentsForMove(moveId);
  }

  toggleShowDocs = () => {
    this.setState({ showDocs: !this.state.showDocs });
  };

  renderHeader = () => {
    const { expenseDocs, weightTicketDocs, inReviewPage } = this.props;
    const totalDocs = expenseDocs.length + weightTicketDocs.length;
    const documentLabel = `document${totalDocs > 1 ? 's' : ''}`;

    return <h3>{inReviewPage ? `Document Summary - ${totalDocs} total` : `${totalDocs} ${documentLabel} added`}</h3>;
  };

  render() {
    const { showDocs } = this.state;
    const { expenseDocs, weightTicketDocs, moveId, showLinks, showToggleDocs, inReviewPage } = this.props;
    const totalDocs = expenseDocs.length + weightTicketDocs.length;

    if (totalDocs === 0) {
      return null;
    }
    return (
      <>
        <div
          className="doc-summary-container"
          data-cy="documents-uploaded"
          style={{ paddingBottom: !inReviewPage && !showDocs ? '1em' : 0, marginTop: !inReviewPage ? '1em' : 0 }}
        >
          <div style={{ display: 'flex', alignItems: 'baseline' }}>
            {this.renderHeader()}
            {showToggleDocs && (
              <a style={{ paddingLeft: '1em' }} onClick={this.toggleShowDocs}>
                {showDocs ? 'Hide' : 'Show'}
              </a>
            )}
          </div>
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
                {formatExpenseDocs(expenseDocs).map(expense => <ExpenseTicketListItem key={expense.id} {...expense} />)}
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
      </>
    );
  }
}

DocumentsUploaded.propTypes = {
  showLinks: bool,
  showToggleDocs: bool,
  inReviewPage: bool,
};

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
