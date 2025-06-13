package common_service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/anglesson/simple-web-server/internal/config"
)

type HubDevService struct {
}

type ConsultaData struct {
	NumeroDeCPF            string `json:"numero_de_cpf"`
	NomeDaPF               string `json:"nome_da_pf"`
	DataNascimento         string `json:"data_nascimento"`
	SituacaoCadastral      string `json:"situacao_cadastral"`
	DataInscricao          string `json:"data_inscricao"`
	DigitoVerificador      string `json:"digito_verificador"`
	ComprovanteEmitido     string `json:"comprovante_emitido"`
	ComprovanteEmitidoData string `json:"comprovante_emitido_data"`
}

type HubDevResponse struct {
	Status   bool         `json:"status"`
	Return   string       `json:"return"`
	Consumed int          `json:"consumed"`
	Result   ConsultaData `json:"result"`
}

func NewHubDevService() *HubDevService {
	return &HubDevService{}
}

func (rf *HubDevService) ConsultaCPF(cpf, dataNascimento string) *HubDevResponse {
	uri := fmt.Sprintf("%s/v2/cpf/?cpf=%s&data=%s&token=%s", config.AppConfig.HubDesenvolvedorApi, cpf, dataNascimento, config.AppConfig.HubDesenvolvedorToken)
	request, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		log.Printf("Erro ao consultar dados na receita federal. Error: %s", err.Error())
		return nil
	}

	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Erro ao fazer requisição para receita federal. Error: %s", err.Error())
		return nil
	}
	defer resp.Body.Close()

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Erro ao ler resposta da receita federal. Error: %s", err.Error())
		return nil
	}

	// Log the raw response
	log.Printf("Resposta bruta da API: %s", string(bodyBytes))

	// Criar um map para armazenar a resposta
	var responseMap map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseMap); err != nil {
		log.Printf("Erro ao fazer parse da resposta: %s", err.Error())
		return nil
	}

	// Extrair o resultado
	resultMap, ok := responseMap["result"].(map[string]interface{})
	if !ok {
		log.Printf("Erro ao extrair resultado da resposta")
		return nil
	}

	// Criar e popular o objeto
	response := &HubDevResponse{
		Status:   responseMap["status"].(bool),
		Return:   responseMap["return"].(string),
		Consumed: int(responseMap["consumed"].(float64)),
		Result: ConsultaData{
			NumeroDeCPF:            resultMap["numero_de_cpf"].(string),
			NomeDaPF:               resultMap["nome_da_pf"].(string),
			DataNascimento:         resultMap["data_nascimento"].(string),
			SituacaoCadastral:      resultMap["situacao_cadastral"].(string),
			DataInscricao:          resultMap["data_inscricao"].(string),
			DigitoVerificador:      resultMap["digito_verificador"].(string),
			ComprovanteEmitido:     resultMap["comprovante_emitido"].(string),
			ComprovanteEmitidoData: resultMap["comprovante_emitido_data"].(string),
		},
	}

	log.Printf("Objeto populado: %+v", response)

	return response
}
