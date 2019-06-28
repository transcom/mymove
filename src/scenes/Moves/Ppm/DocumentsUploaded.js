import React, { Component } from 'react';
import Alert from 'shared/Alert';
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
    const { expenseDocs, weightTicketDocs } = this.props;
    const totalDocs = expenseDocs.length + weightTicketDocs.length;
    if (totalDocs === 0) {
      return null;
    }
    return (
      <>
        {
          <div className="usa-grid" data-cy="documents-uploaded">
            <Alert type="success" heading={this.createHeaderMessage(totalDocs)}>
              {showDocs && (
                <div>
                  {weightTicketDocs.map((ticket, index) => (
                    <WeightTicketListItem key={ticket.id} num={index} {...ticket} />
                  ))}
                  {formatExpenseDocs(expenseDocs).map(expense => (
                    <ExpenseTicketListItem key={expense.id} {...expense} />
                  ))}
                </div>
              )}
            </Alert>
          </div>
        }
      </>
    );
  }
}

function mapStateToProps(state, ownProps) {
  const moveId = ownProps.moveId;
  return {
    moveId: moveId,
    expenseDocs: selectPPMCloseoutDocumentsForMove(state, moveId, ['EXPENSE']),
    weightTicketDocs: selectPPMCloseoutDocumentsForMove(state, moveId, ['WEIGHT_TICKET_SET']),
  };
}

const mapDispatchToProps = {
  selectPPMCloseoutDocumentsForMove,
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(DocumentsUploaded);
