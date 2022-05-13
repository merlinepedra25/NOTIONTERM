package notionterm

import (
	"context"
	"fmt"

	"github.com/ariary/notionion/pkg/notionion"
	"github.com/jomei/notionapi"
)

//GetButtonBlock: retrieve "button" block (embed blocks)
func GetButtonBlock(children notionapi.Blocks) (button notionapi.EmbedBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeEmbed {
			button = *children[i].(*notionapi.EmbedBlock)
			return button, nil
		}
	}
	err = fmt.Errorf("Failed retrieving \"button\" widget")
	return button, err
}

//RequestButtonBlock: retrieve "button" widget (embed block)
func RequestButtonBlock(client *notionapi.Client, pageid string) (terminal notionapi.EmbedBlock, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return terminal, err
	}
	return GetButtonBlock(children)
}

//GetTerminalBlock: retrieve "terminal" block (code blocks)
func GetTerminalBlock(children notionapi.Blocks) (terminal notionapi.CodeBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeCode {
			terminal = *children[i].(*notionapi.CodeBlock)
			//to do check if terminal is under the button
			return terminal, nil
		}
	}
	err = fmt.Errorf("Failed retrieving \"terminal\" section")
	return terminal, err
}

//RequestTerminalBlock: retrieve "terminal" block (code blocks)
func RequestTerminalBlock(client *notionapi.Client, pageid string) (terminal notionapi.CodeBlock, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return terminal, err
	}
	return GetTerminalBlock(children)
}

//RequestTerminalCodeContent: Obtain the content of code block object under the request heading
func RequestTerminalCodeContent(client *notionapi.Client, pageid string) (terminal string, err error) {

	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return "", err
	}
	return GetTerminalCodeContent(children)
}

//GeTerminalCodeContent: Obtain the content of code block object under the request heading whithout making request
func GetTerminalCodeContent(children notionapi.Blocks) (terminal string, err error) {
	termCode, err := GetTerminalBlock(children)
	if err != nil {
		return "", err
	}
	terminal = termCode.Code.RichText[0].PlainText
	return terminal, err
}

//GetTerminalLastRichText: Obtain the last RichText
func GetTerminalLastRichText(termCode notionapi.CodeBlock) (terminal string, err error) {
	terminal = termCode.Code.RichText[len(termCode.Code.RichText)-1].PlainText
	return terminal, err
}

//UpdateButtonUrl: update url of the button widget
func UpdateButtonUrl(client *notionapi.Client, buttonID notionapi.BlockID, url string) (notionapi.Block, error) {
	//construct code block containing request
	widget := notionapi.EmbedBlock{
		Embed: notionapi.Embed{
			Caption: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: notionapi.Text{
						Content: "",
					},
					Annotations: &notionapi.Annotations{
						Bold:          false,
						Italic:        false,
						Strikethrough: false,
						Underline:     false,
						Code:          false,
						Color:         "",
					},
				},
			},
			URL: url,
		},
	}

	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Embed: &widget.Embed,
	}

	return client.Block.Update(context.Background(), buttonID, updateReq)
}

//UpdateButtonCaption: update caption of the given button widget
func UpdateButtonCaption(client *notionapi.Client, button notionapi.EmbedBlock, caption string) (notionapi.Block, error) {
	//construct code block containing request
	widget := button

	captionRich := notionapi.RichText{
		Type: notionapi.ObjectTypeText,
		Text: notionapi.Text{
			Content: caption,
		},
		Annotations: &notionapi.Annotations{
			Bold:   false,
			Italic: true,
			Code:   true,
			Color:  "green",
		},
	}

	widget.Embed.Caption = []notionapi.RichText{captionRich}
	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Embed: &widget.Embed,
	}

	return client.Block.Update(context.Background(), button.ID, updateReq)
}

//UpdateCodeContent: update code block with content
func UpdateCodeContent(client *notionapi.Client, codeBlockID notionapi.BlockID, content string) (notionapi.Block, error) {
	//construct code block containing request
	code := notionapi.CodeBlock{
		Code: notionapi.Code{
			RichText: []notionapi.RichText{
				{
					Type: notionapi.ObjectTypeText,
					Text: notionapi.Text{
						Content: content,
					},
					Annotations: &notionapi.Annotations{
						Bold:          false,
						Italic:        false,
						Strikethrough: false,
						Underline:     false,
						Code:          false,
						Color:         "",
					},
				},
			},
			Language: "shell",
		},
	}

	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Code: &code.Code,
	}

	return client.Block.Update(context.Background(), codeBlockID, updateReq)
}

//AddRichText: Add rich text in code
func AddRichText(client *notionapi.Client, codeBlock notionapi.CodeBlock, content string) (notionapi.Block, error) {
	rich := codeBlock.Code.RichText
	newLine := notionapi.RichText{
		Type: notionapi.ObjectTypeText,
		Text: notionapi.Text{
			Content: content,
		},
	}
	nRich := append(rich, newLine)
	//construct code block containing request
	code := notionapi.CodeBlock{
		Code: notionapi.Code{
			RichText: nRich,
			Language: "shell",
		},
	}
	// send update request
	updateReq := &notionapi.BlockUpdateRequest{
		Code: &code.Code,
	}

	return client.Block.Update(context.Background(), codeBlock.ID, updateReq)
}

//AddTermLine: Add rich text with a new line and "$"
func AddTermLine(client *notionapi.Client, codeBlock notionapi.CodeBlock) (notionapi.Block, error) {

	return AddRichText(client, codeBlock, "$")
}

//GetTableBlock: retrieve table block
func GetTableBlock(children notionapi.Blocks) (table notionapi.TableBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeTableBlock {
			table = *children[i].(*notionapi.TableBlock)
			return table, nil
		}
	}
	err = fmt.Errorf("failed retrieving table block")
	return table, err
}

//RequestTableBlock: retrieve table block by requetsing it
func RequestTableBlock(client *notionapi.Client, pageid string) (table notionapi.TableBlock, err error) {
	children, err := notionion.RequestProxyPageChildren(client, pageid)
	if err != nil {
		return table, err
	}
	return GetTableBlock(children)
}

//GetTableRowBlock: retrieve table row block
func GetTableRowBlock(children notionapi.Blocks) (tableRow notionapi.TableRowBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeTableRowBlock {
			tableRow = *children[i].(*notionapi.TableRowBlock)
			return tableRow, nil
		}
	}
	err = fmt.Errorf("failed retrieving table row block")
	return tableRow, err
}

//RequestTableRowBlock: retrieve table row block by requetsing it
func RequestTableRowBlock(client *notionapi.Client, pageid string) (tableRow notionapi.TableRowBlock, err error) {

	tableBlock, err := RequestTableBlock(client, pageid)
	if err != nil {
		return tableRow, err
	}

	tableBlockChildren, err := client.Block.GetChildren(context.Background(), tableBlock.ID, nil)
	if err != nil {
		return tableRow, err
	}

	return GetTableRowBlock(tableBlockChildren.Results)
}

//GetTableRowBlockbyHeader: retrieve table row block providing its header value
func GetTableRowBlockbyHeader(children notionapi.Blocks, header string) (tableRow notionapi.TableRowBlock, err error) {
	for i := 0; i < len(children); i++ {
		if children[i].GetType() == notionapi.BlockTypeTableRowBlock {
			tableRowTmp := *children[i].(*notionapi.TableRowBlock)
			if len(tableRowTmp.TableRow.Cells) < 0 {
				continue
			}
			if tableRowTmp.TableRow.Cells[0][0].Text.Content == header {
				return tableRowTmp, nil
			}
		}
	}
	err = fmt.Errorf("Failed retrieving table row block")
	return tableRow, err
}

//RequestTableRowBlock: retrieve table row block providing its header valueby requetsing it
func RequestTableRowBlockByHeader(client *notionapi.Client, pageid string, header string) (tableRow notionapi.TableRowBlock, err error) {

	tableBlock, err := RequestTableBlock(client, pageid)
	if err != nil {
		return tableRow, err
	}

	tableBlockChildren, err := client.Block.GetChildren(context.Background(), tableBlock.ID, nil)
	if err != nil {
		return tableRow, err
	}

	return GetTableRowBlockbyHeader(tableBlockChildren.Results, header)
}

func RequestRowValueByHeader(client *notionapi.Client, pageid string, header string) (result string, err error) {
	tableRow, err := RequestTableRowBlockByHeader(client, pageid, header)
	if err != nil {
		return "", err
	}
	if len(tableRow.TableRow.Cells) < 2 {
		err = fmt.Errorf("failed retrieving value in table row (seems that the row does not have more than 1 columns)")
		return "", err
	} else if len(tableRow.TableRow.Cells[1]) < 1 {
		err = fmt.Errorf("failed retrieving value in table row (seems that the value is empty)")
		return "", err
	}
	result = tableRow.TableRow.Cells[1][0].Text.Content
	return result, err
}

func RequestTargetUrl(client *notionapi.Client, pageid string) (targetUrl string, err error) {

	// tableRow, err := RequestTableRowBlock(client, pageid)
	// if err != nil {
	// 	return "", err
	// }
	// // fmt.Printf("%+v", tableRow)
	// t := tableRow.TableRow.Cells[1][0]
	// fmt.Println(t.Text.Content)

	// // for i := 0; i < len(tableRow.TableRow.Cells); i++ {
	// // 	fmt.Println(i)
	// // 	fmt.Printf("%+v", tableRow.TableRow.Cells[i])
	// // }
	return RequestRowValueByHeader(client, pageid, "Target")
}
