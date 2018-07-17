import { get } from 'lodash';
import React, { Component } from 'react';
import { connect } from 'react-redux';
import { bindActionCreators } from 'redux';

import * as Action from 'shared/Entities/actions';
import * as Selector from 'shared/Entities/selectors';
import DocumentContent from './DocumentContent';

export class DocumentUploadViewer extends Component {
  componentDidMount() {
    if (!this.props.moveDocument) {
      this.props.getMoveDocumentsForMove(this.props.match.params.moveId);
    }
  }

  render() {
    let uploadModels = get(this.props.moveDocument, 'document.uploads', []);
    let uploads;
    if (uploadModels.length) {
      uploads = uploadModels.map(upload => (
        <DocumentContent
          key={upload.url}
          url={upload.url}
          filename={upload.filename}
          contentType={upload.content_type}
        />
      ));
    }
    return <div className="document-contents">{uploads}</div>;
  }
}

DocumentUploadViewer.propTypes = {};

function mapStateToProps(state, props) {
  const moveDocumentId = props.match.params.moveDocumentId;
  return {
    entities: state.entities,
    moveDocument: Selector.selectMoveDocument(state, moveDocumentId),
  };
}

function mapDispatchToProps(dispatch) {
  return bindActionCreators(
    {
      getMoveDocumentsForMove: Action.getMoveDocumentsForMove,
    },
    dispatch,
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(
  DocumentUploadViewer,
);
