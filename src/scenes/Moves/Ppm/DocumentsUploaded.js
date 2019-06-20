import React, { Component } from 'react';
import Alert from 'shared/Alert';
import { selectPPMCloseoutDocumentsForMove } from 'shared/Entities/modules/movingExpenseDocuments';
import { getMoveDocumentsForMove } from 'shared/Entities/modules/moveDocuments';
import { connect } from 'react-redux';

export class DocumentsUploaded extends Component {
  componentDidMount() {
    const { moveId } = this.props;
    this.props.getMoveDocumentsForMove(moveId);
  }

  createHeaderMessage = documentLength => {
    return (
      <div>
        {documentLength} document{documentLength > 1 ? 's' : ''} added <a style={{ paddingLeft: '1em' }}>Show</a>
      </div>
    );
  };

  render() {
    const { allDocuments } = this.props;
    const documentLength = allDocuments.length;
    if (documentLength === 0) {
      return null;
    }
    return (
      <>
        {
          <div className="usa-grid" data-cy="documents-uploaded">
            <div className="usa-width-one-whole">
              <Alert type="success" heading={this.createHeaderMessage(documentLength)} />
            </div>
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
    allDocuments: selectPPMCloseoutDocumentsForMove(state, moveId),
  };
}

const mapDispatchToProps = {
  selectPPMCloseoutDocumentsForMove,
  getMoveDocumentsForMove,
};

export default connect(mapStateToProps, mapDispatchToProps)(DocumentsUploaded);
