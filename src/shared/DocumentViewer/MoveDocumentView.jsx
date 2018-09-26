import React, { Component } from 'react';
import PropTypes from 'prop-types';
import DocumentContent from './DocumentContent';
import DocumentList from './DocumentList';
import DocumentDetailPanel from './DocumentDetailPanel';
import { PanelField } from 'shared/EditablePanel';
import { Tab, Tabs, TabList, TabPanel } from 'react-tabs';
import { Link } from 'react-router-dom';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faPlusCircle from '@fortawesome/fontawesome-free-solid/faPlusCircle';

class MoveDocumentView extends Component {
  componentDidMount() {
    const { onDidMount } = this.props;
    onDidMount();
  }

  render() {
    const {
      documentDetailUrlPrefix,
      moveDocuments,
      moveLocator,
      newDocumentUrl,
      serviceMember: { edipi, name },
      uploads,
    } = this.props;
    return (
      <div className="usa-grid doc-viewer">
        <div className="usa-width-two-thirds">
          <div className="tab-content">
            <div className="document-contents">
              {uploads.map(({ url, filename, content_type }) => (
                <DocumentContent
                  key={url}
                  url={url}
                  filename={filename}
                  contentType={content_type}
                />
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
                <Tab className="title nav-tab">
                  All Documents ({moveDocuments.length})
                </Tab>
                <Tab className="title nav-tab">Details</Tab>
              </TabList>

              <TabPanel>
                <div className="pad-ns">
                  <span className="status">
                    <FontAwesomeIcon
                      className="icon link-blue"
                      icon={faPlusCircle}
                    />
                  </span>
                  <Link to={newDocumentUrl}>Upload new document</Link>
                </div>
                <div>
                  {' '}
                  <DocumentList
                    detailUrlPrefix={documentDetailUrlPrefix}
                    moveDocuments={moveDocuments}
                  />
                </div>
              </TabPanel>

              <TabPanel>
                <DocumentDetailPanel
                  className="document-viewer"
                  moveDocumentId={moveDocumentId}
                  title=""
                />
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
