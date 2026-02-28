package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func login(accessKeyId, secretAccessKey string) (token string, err error) {
	url := fmt.Sprintf("%s/authenticate", BASE_URL)
	body := &AuthenticateRequest{
		AccessKeyId:     accessKeyId,
		SecretAccessKey: secretAccessKey,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return token, err
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Println(err)
		return token, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return token, err
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return token, err
	}

	var auth AuthenticateResponse

	err = json.Unmarshal(resp, &auth)
	if err != nil {
		fmt.Println(err)
		return token, err
	}

	return auth.Token, nil
}

func getRuntimeContainerProfile(token string) (profiles []ContainerProfile, err error) {
	const limit = 100
	offset := 0
	allProfiles := []ContainerProfile{}

	for {
		url := fmt.Sprintf("%s/profiles/container", BASE_URL)

		client := &http.Client{}
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Printf("Error creating request at offset %d: %v\n", offset, err)
			return allProfiles, nil
		}

		q := req.URL.Query()
		q.Add("state", "active")
		q.Add("limit", fmt.Sprintf("%d", limit))
		q.Add("offset", fmt.Sprintf("%d", offset))
		req.URL.RawQuery = q.Encode()

		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("Error fetching profiles at offset %d: %v\n", offset, err)
			return allProfiles, nil
		}

		resp, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			fmt.Printf("Error reading response at offset %d: %v\n", offset, err)
			return allProfiles, nil
		}

		var batchProfiles []ContainerProfile
		err = json.Unmarshal(resp, &batchProfiles)
		if err != nil {
			fmt.Printf("Error unmarshaling response at offset %d: %v\n", offset, err)
			return allProfiles, nil
		}

		// Add to all profiles
		allProfiles = append(allProfiles, batchProfiles...)

		fmt.Printf("Fetched %d profiles (total: %d)\n", len(batchProfiles), len(allProfiles))

		// Stop if we got fewer items than the limit
		if len(batchProfiles) < limit {
			fmt.Println("Reached end of profiles")
			break
		}

		// Move to next batch
		offset += limit
	}

	return allProfiles, nil
}

func updateRuntimeContainerPolicy(token string, policy ContainerPolicy) error {
	url := fmt.Sprintf("%s/policies/runtime/container", BASE_URL)

	jsonBody, err := json.Marshal(policy)
	if err != nil {
		return err
	}

	client := &http.Client{}
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("x-prisma-cloud-target-env", `{"permission":"policyRuntimeContainer"}`)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error updating runtime container policy: %v\n", err)
		return err
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return err
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("failed to update policy: %s", string(resp))
	}

	fmt.Printf("Successfully updated runtime container policy\n")
	return nil
}

func getAllRuntimeContainerPolicies(token string) (ContainerPolicy, error) {
	var policy ContainerPolicy

	url := fmt.Sprintf("%s/policies/runtime/container", BASE_URL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return policy, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("x-prisma-cloud-target-env", `{"permission":"policyRuntimeContainer"}`)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching runtime container policies: %v\n", err)
		return policy, err
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return policy, err
	}

	err = json.Unmarshal(resp, &policy)
	if err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
		return policy, err
	}

	fmt.Printf("Successfully fetched runtime container policy with %d rules\n", len(policy.Rules))
	return policy, nil
}

func getAllRuntimeHostPolicies(token string) (HostPolicy, error) {
	var policy HostPolicy

	url := fmt.Sprintf("%s/policies/runtime/host", BASE_URL)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return policy, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("x-prisma-cloud-target-env", `{"permission":"policyRuntimeHosts"}`)

	res, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error fetching runtime host policies: %v\n", err)
		return policy, err
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return policy, err
	}

	err = json.Unmarshal(resp, &policy)
	if err != nil {
		fmt.Printf("Error unmarshaling response: %v\n", err)
		return policy, err
	}

	fmt.Printf("Successfully fetched runtime host policy with %d rules\n", len(policy.Rules))
	return policy, nil
}
