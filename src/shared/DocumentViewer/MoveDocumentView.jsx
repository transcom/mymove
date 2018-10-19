import React, { Component } from 'react';
import PropTypes from 'prop-types';
import DocumentContent from './DocumentContent';
import DocumentList from './DocumentList';
import { PanelField } from 'shared/EditablePanel';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';
import './index.css';
import DocumentDetailPanelView from 'shared/DocumentViewer/DocumentDetailPanelView';

class MoveDocumentView extends Component {
  componentDidMount() {
    const { onDidMount } = this.props;
    onDidMount();
  }

  render() {
    const {
      documentDetailUrlPrefix,
      moveDocument,
      moveDocumentSchema,
      moveDocuments,
      moveLocator,
      newDocumentUrl,
      serviceMember: { edipi, name },
      uploads,
    } = this.props;
    const currentMoveDocumentId = this.props.match.params.moveDocumentId;
    return (
      <div className="usa-grid doc-viewer">
        <div className="usa-width-two-thirds">
          <div className="tab-content">
            <div className="document-contents">
              {uploads.map(({ url, filename, content_type }) => (
                <DocumentContent key={url} url={url} filename={filename} contentType={content_type} />
              ))}
            </div>
          </div>
        </div>
        <div className="usa-width-one-third">
          <h3>{name}</h3>
          <PanelField title="Move Locator">{moveLocator}</PanelField>
          <PanelField title="DoD ID">{edipi}</PanelField>
          <div className="tab-content">
            <Tabs defaultIndex={0}>
              <TabList className="doc-viewer-tabs">
                <Tab className="title nav-tab">All Documents ({moveDocuments.length})</Tab>
                <Tab className="title nav-tab">Details</Tab>
              </TabList>

              <TabPanel>
                <div className="pad-ns">
                  <span className="status">
                    <FontAwesomeIcon className="icon link-blue" icon={faPlusCircle} />
                  </span>
                  <Link to={newDocumentUrl}>Upload new document</Link>
                </div>
                <div>
                  {' '}
                  <DocumentList
                    currentMoveDocumentId={currentMoveDocumentId}
                    detailUrlPrefix={documentDetailUrlPrefix}
                    moveDocuments={moveDocuments}
                  />
                </div>
              </TabPanel>

              <TabPanel>
                <DocumentDetailPanelView schema={moveDocumentSchema} {...moveDocument} />
              </TabPanel>
            </Tabs>
          </div>
        </div>
      </div>
    );
  }
}

MoveDocumentView.propTypes = {
  documentDetailUrlPrefix: PropTypes.string.isRequired,
  moveDocument: PropTypes.shape({
    createdAt: PropTypes.string.isRequired,
    notes: PropTypes.string.isRequired,
    status: PropTypes.string.isRequired,
    title: PropTypes.string.isRequired,
    type: PropTypes.string.isRequired,
  }),
  moveDocumentSchema: PropTypes.object.isRequired,
  moveDocuments: PropTypes.array.isRequired,
  moveLocator: PropTypes.string.isRequired,
  newDocumentUrl: PropTypes.string.isRequired,
  onDidMount: PropTypes.func.isRequired,
  serviceMember: PropTypes.shape({
    edipi: PropTypes.string.isRequired,
    name: PropTypes.string.isRequired,
  }),
  uploads: PropTypes.array.isRequired,
};

export default MoveDocumentView;
