import React from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { faFile, faPen, faPlus, faChevronLeft, faChevronRight, faChevronDown, faChevronUp, faCheck, faTimes, faSearchPlus, faSearchMinus, faRedoAlt, faUndoAlt, faLock, faMapMarkerAlt, faArrowRight, faArrowLeft, faThList, faQuestionCircle, faPhoneAlt, faClock, faPlusCircle, faPlayCircle, faPlusSquare} from '@fortawesome/free-solid-svg-icons';
import { faTimesCircle, faCalendar } from '@fortawesome/free-regular-svg-icons';
import { faCheckCircle as fasCheckCircle } from '@fortawesome/free-solid-svg-icons'
import { faCheckCircle as farCheckCircle } from '@fortawesome/free-regular-svg-icons';

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
        <FontAwesomeIcon icon={faPen} />
        <code>edit | faPen</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faPlus} />
        <code>add | faPlus</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faChevronLeft} />
        <code>chevron-left |  faChevronLeft</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faChevronRight} />
        <code>chevron-right |  faChevronRight</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faChevronDown} />
        <code>chevron-down |  faChevronDown</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faChevronUp} />
        <code>chevron-up |  faChevronUp</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faCheck} />
        <code>checkmark | faCheck</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faTimes} />
        <code> x |  faTimes</code>
      </div>
      <div>
        <FontAwesomeIcon icon={farCheckCircle} />
        <code>accept | faCheckCircle (regular)</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faTimesCircle} />
        <code>reject | faTimesCircle (regular)</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faSearchPlus} />
        <code>zoom in | faSearchPlus</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faSearchMinus} />
        <code>zoom out | faSearchMinus</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faRedoAlt} />
        <code>rotate clockwise | faRedoAlt</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faUndoAlt} />
        <code>rotate counter clockwise | faUndoAlt</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faLock} />
        <code>lock | faLock</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faMapMarkerAlt} />
        <code>map pin | faMapMarkerAlt</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faArrowRight} />
        <code>arrow right | faArrowRight</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faArrowLeft} />
        <code>arrow left| faArrowLeft</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faThList} />
        <code>doc menu | faThList</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faCalendar} />
        <code>calendar | faCalendar</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faQuestionCircle} />
        <code>question circle | faQuestionCircle</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faPhoneAlt} />
        <code>phone | faPhoneAlt</code>
      </div>
      <div>
        <FontAwesomeIcon icon={fasCheckCircle} />
        <code>check circle | faCheckCircle (solid)</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faClock} />
        <code>clock | faClock</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faPlusCircle} />
        <code>plus circle | faPlusCircle</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faPlayCircle} />
        <code>play circle | faPlayCircle</code>
      </div>
      <div>
        <FontAwesomeIcon icon={faPlusSquare} />
        <code>plus square | faPlusSquare</code>
      </div>
    </div>
  </div>
);
