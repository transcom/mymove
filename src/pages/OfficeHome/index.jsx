import { connect } from 'react-redux';
import OfficeHome from './OfficeHome';

const mapStateToProps = (state) => {
  return {
    activeRole: state.auth.activeRole,
  };
};

export default connect(mapStateToProps)(OfficeHome);
