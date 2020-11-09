import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFile, faPen, faPlus, faChevronLeft, faChevronRight, faChevronDown, faChevronUp, faCheck, faTimes, faSearchPlus, faSearchMinus, faRedoAlt, faUndoAlt, faLock, faSolidMapMarkerAlt, faArrowRight, faArrowLeft, faThList} from '@fortawesome/free-solid-svg-icons';
import { faCheckCircle, faTimesCircle, faCalendar } from '@fortawesome/free-regular-svg-icons';

// Icons
export default {
  title: 'Global|Icons',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/eabef4a2-603e-4c3d-b249-0e580f1c8306?mode=design',
    },
  },
};

export const all = () => (
  <div style={{ padding: '20px', background: '#f0f0f0' }}>
    <h3>Icons</h3>
    <div id="icons" style={{ display: 'flex', flexWrap: 'wrap' }}>
      <div>
        <FontAwesomeIcon icon={faFile} />
        <code>documents | faFile</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faCheckCircle} />
        <code>accept | faCheckCircle (regular)</code>
      </div>
    </div>
  </div>
);
