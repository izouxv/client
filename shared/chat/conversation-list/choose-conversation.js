// @flow
import * as React from 'react'
import * as Types from '../../constants/types/chat2'
import * as Kb from '../../common-adapters'
import * as Styles from '../../styles'
import ConversationList from './conversation-list-container'

type Props = {|
  dropdownButtonStyle?: ?Styles.StylesCrossPlatform,
  overlayStyle?: ?Styles.StylesCrossPlatform,
  filter?: string,
  onSelect: (conversationIDKey: Types.ConversationIDKey) => void,
  onSetFilter?: (filter: string) => void,
  selectedText: string, // TODO: make this pretty
|}

type State = {|expanded: boolean|}

class ChooseConversation extends React.Component<Props & Kb.OverlayParentProps, State> {
  state = {
    expanded: false,
  }
  _toggleOpen = () => {
    this.setState(prevState => ({expanded: !prevState.expanded}))
  }
  _onSelect = (conversationIDKey: Types.ConversationIDKey) => {
    this.setState({expanded: false})
    this.props.onSelect(conversationIDKey)
  }
  render() {
    return (
      <>
        <Kb.DropdownButton
          disabled={false}
          selected={<Kb.Text type="BodySemibold">{this.props.selectedText}</Kb.Text>}
          setAttachmentRef={this.props.setAttachmentRef}
          toggleOpen={this._toggleOpen}
          style={Styles.collapseStyles([styles.dropdownButton, this.props.dropdownButtonStyle])}
        />
        <Kb.Overlay
          attachTo={this.props.getAttachmentRef}
          onHidden={this._toggleOpen}
          position={'center center'}
          style={Styles.collapseStyles([styles.overlay, this.props.overlayStyle])}
          visible={this.state.expanded}
        >
          <ConversationList
            onSelect={this._onSelect}
            filter={this.props.filter}
            onSetFilter={this.props.onSetFilter}
          />
        </Kb.Overlay>
      </>
    )
  }
}

export default Kb.OverlayParentHOC(ChooseConversation)

const styles = Styles.styleSheetCreate({
  dropdownButton: {
    width: 300,
  },
  overlay: {
    backgroundColor: Styles.globalColors.white,
    height: 360,
    width: 300,
  },
})
